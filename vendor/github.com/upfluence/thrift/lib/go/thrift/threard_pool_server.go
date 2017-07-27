/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package thrift

import (
	"log"
	"runtime/debug"
)

// ThreardPool, non-concurrent server for testing.
type TThreardPoolServer struct {
	quit        chan struct{}
	requestChan chan int

	processorFactory       TProcessorFactory
	serverTransport        TServerTransport
	inputTransportFactory  TTransportFactory
	outputTransportFactory TTransportFactory
	inputProtocolFactory   TProtocolFactory
	outputProtocolFactory  TProtocolFactory
}

func NewTThreardPoolServer2(processor TProcessor, serverTransport TServerTransport, poolSize int) *TThreardPoolServer {
	return NewTThreardPoolServerFactory2(
		NewTProcessorFactory(processor),
		serverTransport,
		poolSize,
	)
}

func NewTThreardPoolServer4(processor TProcessor, serverTransport TServerTransport, transportFactory TTransportFactory, protocolFactory TProtocolFactory, poolSize int) *TThreardPoolServer {
	return NewTThreardPoolServerFactory4(
		NewTProcessorFactory(processor),
		serverTransport,
		transportFactory,
		protocolFactory,
		poolSize,
	)
}

func NewTThreardPoolServer6(processor TProcessor, serverTransport TServerTransport, inputTransportFactory TTransportFactory, outputTransportFactory TTransportFactory, inputProtocolFactory TProtocolFactory, outputProtocolFactory TProtocolFactory, poolSize int) *TThreardPoolServer {
	return NewTThreardPoolServerFactory6(
		NewTProcessorFactory(processor),
		serverTransport,
		inputTransportFactory,
		outputTransportFactory,
		inputProtocolFactory,
		outputProtocolFactory,
		poolSize,
	)
}

func NewTThreardPoolServerFactory2(processorFactory TProcessorFactory, serverTransport TServerTransport, poolSize int) *TThreardPoolServer {
	return NewTThreardPoolServerFactory6(
		processorFactory,
		serverTransport,
		NewTTransportFactory(),
		NewTTransportFactory(),
		NewTBinaryProtocolFactoryDefault(),
		NewTBinaryProtocolFactoryDefault(),
		poolSize,
	)
}

func NewTThreardPoolServerFactory4(processorFactory TProcessorFactory, serverTransport TServerTransport, transportFactory TTransportFactory, protocolFactory TProtocolFactory, poolSize int) *TThreardPoolServer {
	return NewTThreardPoolServerFactory6(
		processorFactory,
		serverTransport,
		transportFactory,
		transportFactory,
		protocolFactory,
		protocolFactory,
		poolSize,
	)
}

func NewTThreardPoolServerFactory6(processorFactory TProcessorFactory, serverTransport TServerTransport, inputTransportFactory TTransportFactory, outputTransportFactory TTransportFactory, inputProtocolFactory TProtocolFactory, outputProtocolFactory TProtocolFactory, poolSize int) *TThreardPoolServer {
	return &TThreardPoolServer{
		processorFactory:       processorFactory,
		serverTransport:        serverTransport,
		inputTransportFactory:  inputTransportFactory,
		outputTransportFactory: outputTransportFactory,
		inputProtocolFactory:   inputProtocolFactory,
		outputProtocolFactory:  outputProtocolFactory,
		quit:        make(chan struct{}, 1),
		requestChan: make(chan int, poolSize),
	}
}

func (p *TThreardPoolServer) ProcessorFactory() TProcessorFactory {
	return p.processorFactory
}

func (p *TThreardPoolServer) ServerTransport() TServerTransport {
	return p.serverTransport
}

func (p *TThreardPoolServer) InputTransportFactory() TTransportFactory {
	return p.inputTransportFactory
}

func (p *TThreardPoolServer) OutputTransportFactory() TTransportFactory {
	return p.outputTransportFactory
}

func (p *TThreardPoolServer) InputProtocolFactory() TProtocolFactory {
	return p.inputProtocolFactory
}

func (p *TThreardPoolServer) OutputProtocolFactory() TProtocolFactory {
	return p.outputProtocolFactory
}

func (p *TThreardPoolServer) Listen() error {
	return p.serverTransport.Listen()
}

func (p *TThreardPoolServer) AcceptLoop() error {
	for {
		client, err := p.serverTransport.Accept()
		if err != nil {
			select {
			case <-p.quit:
				return nil
			default:
			}
			return err
		}
		if client != nil {
			p.requestChan <- 1
			go func() {
				if err := p.processRequests(client); err != nil {
					log.Println("error processing request:", err)
				}

				<-p.requestChan
			}()
		}
	}
}

func (p *TThreardPoolServer) Serve() error {
	err := p.Listen()
	if err != nil {
		return err
	}
	p.AcceptLoop()
	return nil
}

func (p *TThreardPoolServer) Stop() error {
	p.quit <- struct{}{}
	p.serverTransport.Interrupt()
	return nil
}

func (p *TThreardPoolServer) processRequests(client TTransport) error {
	processor := p.processorFactory.GetProcessor(client)
	inputTransport := p.inputTransportFactory.GetTransport(client)
	outputTransport := p.outputTransportFactory.GetTransport(client)
	inputProtocol := p.inputProtocolFactory.GetProtocol(inputTransport)
	outputProtocol := p.outputProtocolFactory.GetProtocol(outputTransport)
	defer func() {
		if e := recover(); e != nil {
			log.Printf("panic in processor: %s: %s", e, debug.Stack())
		}
	}()
	if inputTransport != nil {
		defer inputTransport.Close()
	}
	if outputTransport != nil {
		defer outputTransport.Close()
	}
	for {
		ok, err := processor.Process(inputProtocol, outputProtocol)
		if err, ok := err.(TTransportException); ok && err.TypeId() == END_OF_FILE {
			return nil
		} else if err != nil {
			log.Printf("error processing request: %s", err)
			return err
		}
		if !ok {
			break
		}
	}
	return nil
}

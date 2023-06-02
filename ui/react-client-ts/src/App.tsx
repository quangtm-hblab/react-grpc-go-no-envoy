import React from "react";
import logo from "./logo.svg";
import "./App.css";
import * as grpcWeb from "grpc-web";
import { SumRequest, SumResponse } from "./protobuf/calculator_pb";
import { CalculateClient } from "./protobuf/CalculatorServiceClientPb";
function App() {
  const callSum = () => {
    let calClient: CalculateClient = new CalculateClient(
      "https://localhost:50022"
    );
    const req: SumRequest = new SumRequest();
    req.setNum1(1);
    req.setNum2(2);
    calClient.sum(req, {}, (error: grpcWeb.RpcError, res: SumResponse) => {
      console.log(res.getResult());
    });
  };

  return (
    <div>
      <button onClick={callSum}>Click</button>
    </div>
  );
}

export default App;

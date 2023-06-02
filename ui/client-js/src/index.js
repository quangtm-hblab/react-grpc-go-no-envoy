const { SumRequest, SumResponse } = require("./protobuf/calculator_pb.js");
const { CalculateClient } = require("./protobuf/calculator_grpc_web_pb.js");

var client = new CalculateClient("https://localhost:50022");
var request = new SumRequest();
request.setNum1(1);
request.setNum2(2);

client.sum(request, {}, function (err, res) {
  console.log(res);
});

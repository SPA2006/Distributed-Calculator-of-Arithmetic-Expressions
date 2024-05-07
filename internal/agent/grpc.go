package grpc_agent

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	comp "github.com/SPA2006/Distributed-Calculator-of-Arithmetic-Expressions/internal/agent"
	database "github.com/SPA2006/Distributed-Calculator-of-Arithmetic-Expressions/internal/database"
	server "github.com/SPA2006/Distributed-Calculator-of-Arithmetic-Expressions/internal/server/businessLogic/rpn"
	pb "github.com/SPA2006/Distributed-Calculator-of-Arithmetic-Expressions/proto/agent/agent"
)

type Expression struct {
	ID                 int64  `json:"id"`
	Expression         string `json:"expression"`
	UserID             int64  `json:"user_id"`
	TimePower          int64  `json:"time_power"`
	TimeMultiplication int64  `json:"time_mult"`
	TimeDivision       int64  `json:"time_div"`
	TimeAddition       int64  `json:"time_add"`
	TimeSubstraction   int64  `json:"time_sub"`
	Result             int64  `json:"result"`
}

type ServerAPI struct {
	pb.UnimplementedAgentServer
	expression Expression
}

func NewServer() *ServerAPI {
	return &ServerAPI{}
}

func (s *ServerAPI) CalculateExpression(ctx context.Context,
	reqExp *pb.ExpressionToCalculateRequest) (*pb.Response, error) {
	// get expression from gRPC request and 'translate' it into RPN
	in := reqExp.GetExpresion()
	rpn_in, err := server.Get_rpn(in)

	if err != nil {
		return nil, err
	}

	power := reqExp.GetTimePower()
	mult := reqExp.GetTimeMultiplication()
	div := reqExp.GetTimeDivision()
	add := reqExp.GetTimeAddition()
	sub := reqExp.GetTimeSubstraction()

	// setting default values
	if power == 0 {
		power = 10
	}
	if mult == 0 {
		mult = 10
	}
	if div == 0 {
		div = 10
	}
	if add == 0 {
		add = 10
	}
	if sub == 0 {
		sub = 10
	}

	// getting the result of calculating (computing) an expression and retruning its result
	result, err := comp.ComputeRPN(rpn_in)
	strResult := strconv.Itoa(result)

	if err != nil {
		return nil, err
	}

	id := reqExp.GetId()
	var dbCopy database.Storage
	err = dbCopy.UpdateExpressionByID(ctx, strResult, id)

	if err != nil {
		return nil, err
	}

	total, err := timeExpression(rpn_in, power, mult, div, add, sub)
	log.Println("waiting for: ", total, " seconds")
	time.Sleep(time.Duration(total) * time.Second)

	log.Println("success (OK) 200")
	resultStr := strconv.Itoa(result)

	return &pb.Response{
		Id:     reqExp.Id,
		Result: resultStr,
	}, nil
}

func timeExpression(rpn_in string, pow, mult, div, add, sub int64) (int64, error) {
	sepExpression := strings.Split(rpn_in, " ")
	var total_time int64
	total_time = 0

	for _, value := range sepExpression {
		switch value {
		case "^":
			total_time += pow
		case "*":
			total_time += mult
		case "/":
			total_time += div
		case "+":
			total_time += add
		case "-":
			total_time += sub
		}
	}

	return total_time, nil
}

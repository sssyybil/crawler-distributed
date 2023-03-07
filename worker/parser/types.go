package parser

import (
	"crawler-distributed/config"
	"crawler-distributed/model"
	"crawler-distributed/pb"
	"crawler-distributed/zhenai/parser"
	"errors"
	"log"
)

/**【序列化】**/

// SerializeRequest 因 engine.Request 中包含函数，无法在网络中传输，所以将其转换为 SerializedRequest 类型。【For Worker Client】
func SerializeRequest(r model.Request) *pb.SerializedRequest {
	name, args := r.Parser.Serialize()
	return &pb.SerializedRequest{
		Url: r.Url,
		Parser: &pb.ParserFunc{
			FunctionName: name,
			Args:         args,
		},
	}
}

// SerializeParserResult 【For Worker Server】
func SerializeParserResult(parserResult model.ParserResult) *pb.SerializedParserResult {
	var reqs []*pb.SerializedRequest
	for _, r := range parserResult.Requests {
		reqs = append(reqs, SerializeRequest(r))
	}

	var items []*pb.Item
	items = modelItemToPbItem(parserResult.Items)
	return &pb.SerializedParserResult{
		Items:   items,
		Request: reqs,
	}
}

/**【反序列化】**/

// DeserializeRequest 【For Worker Client】
func DeserializeRequest(r *pb.SerializedRequest) (model.Request, error) {
	parserFunc, err := getFunctionByNameAndArgs(r.Parser)
	if err != nil {
		return model.Request{}, err
	}

	return model.Request{
		Url:    r.Url,
		Parser: parserFunc,
	}, nil
}

// getFunctionByNameAndArgs 根据 函数名➕参数 获取对应的函数
func getFunctionByNameAndArgs(p *pb.ParserFunc) (model.Parser, error) {
	var tp *model.FuncParser
	switch p.FunctionName {
	case config.ParseCityList:
		tp = model.NewFuncParser(parser.ParseCityList, config.ParseCityList)
	case config.ParseCity:
		tp = model.NewFuncParser(parser.ParseCity, config.ParseCity)
	case config.ParseProfile:
		return parser.NewProfileParser(p.Args), nil
	default:
		return nil, errors.New("[getFunctionByNameAndArgs] unknown parser name")
	}
	return tp, nil
}

// DeserializeParserResult 【For Worker Server】
func DeserializeParserResult(r *pb.SerializedParserResult) model.ParserResult {
	var reqs []model.Request
	for _, r := range r.Request {
		request, err := DeserializeRequest(r)
		if err != nil {
			log.Printf("error deserializing requesting: %v\n", err)
			continue
		}
		reqs = append(reqs, request)
	}
	return model.ParserResult{
		Requests: reqs,
		Items:    pbItemToModelItem(r.Items),
	}
}

func modelItemToPbItem(items []model.Item) []*pb.Item {
	var res []*pb.Item
	for _, v := range items {
		item := &pb.Item{
			Url:  v.Url,
			Id:   v.Id,
			Type: v.Type,
		}
		switch t := v.Payload.(type) {
		case model.Profile:
			p := pb.Profile{
				Name:       t.Name,
				Gender:     t.Gender,
				Age:        uint32(t.Age),
				Height:     float32(t.Height),
				Weight:     float32(t.Weight),
				Income:     t.Income,
				Marriage:   t.Marriage,
				Education:  t.Education,
				Occupation: t.Occupation,
				HuKou:      t.HuKou,
				XinZuo:     t.XinZuo,
				House:      t.House,
				Car:        t.Car,
				CommonInfo: t.CommonInfo,
			}
			item.Profile = &p
		default:
			log.Fatalf("[modelItemToPbItem] convert FAIL! expected type=%T\n", t)
		}
		res = append(res, item)
	}
	return res
}

func pbItemToModelItem(items []*pb.Item) []model.Item {
	var res []model.Item
	for _, v := range items {
		item := model.Item{
			Url:     v.Url,
			Id:      v.Id,
			Type:    v.Type,
			Payload: v.Profile,
		}
		res = append(res, item)
	}
	return res
}

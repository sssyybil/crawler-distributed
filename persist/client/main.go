package client

import (
	"context"
	"crawler-distributed/model"
	"crawler-distributed/pb"
	"crawler-distributed/support/grpcsupport"
	"log"
)

func StartItemSaverClient(address string) (chan model.Item, error) {
	log.Println("💫ItemSaver client is running...")

	grpcClient := grpcsupport.NewItemSaverClient(address)
	ctx := context.Background()
	out := make(chan model.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("[ItemSaver] got item #%d: %v", itemCount, item)
			itemCount++

			pbItem := engineItemToPbItem(item)
			_, err := grpcClient.Save(ctx, &pb.ItemSaverRequest{Item: pbItem})
			if err != nil {
				log.Printf("Item Saver: error saving item %v: %v\n", item, err)
			}
		}
	}()
	return out, nil
}

func engineItemToPbItem(item model.Item) *pb.Item {
	switch p := item.Payload.(type) {
	case *pb.Profile:
		payload := pb.Profile{
			Name:       p.Name,
			Gender:     p.Gender,
			Age:        p.Age,
			Height:     p.Height,
			Weight:     p.Weight,
			Income:     p.Income,
			Marriage:   p.Marriage,
			Education:  p.Education,
			Occupation: p.Occupation,
			HuKou:      p.HuKou,
			XinZuo:     p.XinZuo,
			House:      p.House,
			Car:        p.Car,
			CommonInfo: p.CommonInfo,
		}
		return &pb.Item{
			Url:     item.Url,
			Id:      item.Id,
			Type:    item.Type,
			Profile: &payload,
		}
	default:
		log.Fatalf("[engineItemToPbItem] error to convert model.Item to pb.Item: %T", p)
		return nil
	}
}

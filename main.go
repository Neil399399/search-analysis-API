package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"search-analysis-API/datamodel"
	"search-analysis-API/search"
	"strconv"
)

var (
	port   = "80"
	Search datamodel.Search
)

func main() {
	//http server
	myFunction := func() {
		//handle
		http.HandleFunc("/search", DataSearch)
		http.HandleFunc("/analysis", DataAnalysis)
		http.HandleFunc("/search-analysis", DataSearch_Analysis)
		http.HandleFunc("/search-mock", DataSearch_Mock)

		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			panic("Connect Fail:" + err.Error())
		}
	}
	go myFunction()
	// use go channel to continous code
	endChannel := make(chan os.Signal)
	signal.Notify(endChannel)
	sig := <-endChannel
	fmt.Println("END!:", sig)
}

//function
// /search?APIKEY=AIzaSyCigqPQLr341O-UL_jyJQNdX76fO0TtywA$KEYWORD:=海鮮餐廳&LAT=25.03978&LNG=121.548495
func DataSearch(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	apikey := req.FormValue("APIKEY")
	lat := req.FormValue("LAT")
	lng := req.FormValue("LNG")
	keyword := req.FormValue("KEYWORD")

	//check lat,lng format from http
	if len(lat) != 0 {
		lat64, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			fmt.Println("LAT has wrong format !!!")
			return
		}
		Search.LAT = lat64
	}
	if len(lng) != 0 {
		lng64, err := strconv.ParseFloat(lng, 64)
		if err != nil {
			fmt.Println("LNG has wrong format !!!")
			return
		}
		Search.LNG = lng64
	}

	//check from client
	err := search.Initialize(apikey, 500)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	search500 := search.NewSearch(apikey, 500)
	Search.KEYWORD = keyword
	if !Search.Verify(Search) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "json")

	//search
	List, err := search500.Place(Search.KEYWORD, Search.LAT, Search.LNG)
	if err != nil {
		fmt.Println("google Place Search Error!!", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	//convert to json, give to fprint
	b, err := json.Marshal(List)
	if err != nil {
		fmt.Println("Json Marchal Error!!", err)
	}
	fmt.Fprint(w, string(b))
}

func DataAnalysis(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var top [3]string

	if req.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		type RequestMessage struct {
			Params []string
			Data   []datamodel.Coffee
		}

		var requestMessage RequestMessage
		//check err
		err = json.Unmarshal(body, &requestMessage)
		if err != nil {
			fmt.Println("Json Unmarshal Error!!", err)
			return
		}
		//run jieba
		jiebres, err := Jiebatest(requestMessage.Data, requestMessage.Params)
		if err == ErrNoData {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		if err == ErrIndexSearch {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		//count total
		sortres, err := SortTotal(jiebres)
		if err != nil {
			fmt.Println("Sort Total Error!!", err)
		}
		//find top3
		first, second, third, err := Top3(sortres)
		if err != nil {
			fmt.Println("Find Top3 Error!!", err)
		}
		//print top3
		top1, top2, top3, err := FindIDInfo(first, second, third, requestMessage.Data)
		if err != nil {
			fmt.Println("json marshal failed!!", err)
		}
		top[0] = top1
		top[1] = top2
		top[2] = top3
		b, err := json.Marshal(top)
		fmt.Fprint(w, string(b))
		if err != nil {
			fmt.Println("Find ID Info Error!!", err)
		}

	}
}

func DataSearch_Analysis(w http.ResponseWriter, req *http.Request) {
	var top [3]string
	//set header to tell server which http domain can connect
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	apikey := req.FormValue("APIKEY")
	lat := req.FormValue("LAT")
	lng := req.FormValue("LNG")
	keyword := req.FormValue("KEYWORD")
	name1 := req.FormValue("analysis_word1")
	name2 := req.FormValue("analysis_word2")
	name3 := req.FormValue("analysis_word3")
	querys = []string{name1, name2, name3}

	//check lat,lng format from http
	if len(lat) != 0 {
		lat64, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			fmt.Println("LAT has wrong format !!!")
			return
		}
		Search.LAT = lat64
	}
	if len(lng) != 0 {
		lng64, err := strconv.ParseFloat(lng, 64)
		if err != nil {
			fmt.Println("LNG has wrong format !!!")
			return
		}
		Search.LNG = lng64
	}
	//check from client
	err := search.Initialize(apikey, 500)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	search500 := search.NewSearch(apikey, 500)
	Search.KEYWORD = keyword
	if !Search.Verify(Search) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	//header set
	w.Header().Set("Content-Type", "json")
	//search
	List, err := search500.Place(keyword, Search.LAT, Search.LNG)
	if err != nil {
		fmt.Println("google Place Search Error!!", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	fmt.Println("google  Search Success!!")
	//Analysis
	jiebres, err := Jiebatest(List, querys)
	if err != nil {
		fmt.Println("jieba Error!!", err)
	}
	//count total
	sortres, err := SortTotal(jiebres)
	if err != nil {
		fmt.Println("Sort Total Error!!", err)
	}

	//find top3
	first, second, third, err := Top3(sortres)
	if err != nil {
		fmt.Println("Find Top3 Error!!", err)
	}

	//print top3

	top1, top2, top3, err := FindIDInfo(first, second, third, List)
	if err != nil {
		fmt.Println("json marshal failed!!", err)
	}
	top[0] = top1
	top[1] = top2
	top[2] = top3
	b, err := json.Marshal(top)
	fmt.Fprint(w, string(b))
	if err != nil {
		fmt.Println("Find ID Info Error!!", err)
	}

}

func DataSearch_Mock(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var mockData = `[{"Id":"ChIJSTLZ6barQjQRDNycA51cBq4","Name":"Starbucks 統一星巴克 (101門市)","Rate":4.2,"Addr":"110台灣台北市信義區市府路45號1樓","Reviews":[{"StoreId":"ChIJSTLZ6barQjQRDNycA51cBq4","Text":"不用上101也可以享受風景"},{"StoreId":"ChIJSTLZ6barQjQRDNycA51cBq4","Text":"Make sure you have a reservation prior to going up here. It's a nice view from here. Enjoy the coffee."},{"StoreId":"ChIJSTLZ6barQjQRDNycA51cBq4","Text":"很意外星巴克竟然會藏在101的一樓內，小小的空間但是剛好滿足想喝點東西的需求"},{"StoreId":"ChIJSTLZ6barQjQRDNycA51cBq4","Text":"這裡應該是台北市最高的星巴克⋯!\n或許這裡的租金不便宜吧！比起來並沒有其他分店來的大！只是有著寬廣的視野！也是來這裡第一個原因！\n在這裡消費有分員工與非員工～非員工有最低消費的規定喔！\n但最重要還是看看你要的視野吧！"},{"StoreId":"ChIJSTLZ6barQjQRDNycA51cBq4","Text":"景觀很好，但人很多比較吵ㄧ點！"}]},{"Id":"ChIJiw1lNLqrQjQR7ArhgEV5Sjs","Name":"MR.BROWN 伯朗咖啡館 信義店","Rate":3.7,"Addr":"110台灣台北市信義區松高路12號3樓","Reviews":[{"StoreId":"ChIJiw1lNLqrQjQR7ArhgEV5Sjs","Text":"好喝，份量剛好奶泡滿杯線，可惜9/4要結束營業"},{"StoreId":"ChIJiw1lNLqrQjQR7ArhgEV5Sjs","Text":"拿鐵如果提醒有加糖會更好，無糖一定不錯"},{"StoreId":"ChIJiw1lNLqrQjQR7ArhgEV5Sjs","Text":"奶泡很雷，要好喝奶泡，要遇到好店員，目前買一送一日，奶泡很空虛。"},{"StoreId":"ChIJiw1lNLqrQjQR7ArhgEV5Sjs","Text":"Very bad service attitude and the worst coffee ever - from a costumer lives in a walking distance for more than 20 years. "},{"StoreId":"ChIJiw1lNLqrQjQR7ArhgEV5Sjs","Text":"不錯 只是店有點小"}]},{"Id":"ChIJo_L6IbGrQjQRPlLIE8O5aRs","Name":"考非 Coffee（信義）","Rate":5,"Addr":"11049台灣台北市信義區信義路五段15號B1","Reviews":[{"StoreId":"ChIJo_L6IbGrQjQRPlLIE8O5aRs","Text":"咖啡好喝又便宜!"},{"StoreId":"ChIJo_L6IbGrQjQRPlLIE8O5aRs","Text":""}]},{"Id":"ChIJewqtNbqrQjQRaxR_poZHNPY","Name":"cama café -新光A8店","Rate":3.9,"Addr":"110台灣台北市信義區松高路12號B2","Reviews":[{"StoreId":"ChIJewqtNbqrQjQRaxR_poZHNPY","Text":"平日和假日人潮出現前來較佳，算是cp值高的良品咖啡店，帶本書享受時間"},{"StoreId":"ChIJewqtNbqrQjQRaxR_poZHNPY","Text":"海鹽咖啡超甜，冰黑咖啡超稀薄像喝咖啡水，有史以來最難喝的cama......"},{"StoreId":"ChIJewqtNbqrQjQRaxR_poZHNPY","Text":"還是選擇我的黑咖啡。"},{"StoreId":"ChIJewqtNbqrQjQRaxR_poZHNPY","Text":"喜歡紅茶拿鐵"},{"StoreId":"ChIJewqtNbqrQjQRaxR_poZHNPY","Text":"休閒找咖啡喝的好去處"}]},{"Id":"ChIJP-nV3LmrQjQRfV4-XIaznRA","Name":"Louisa Coffee 路易．莎咖啡(北市府店)","Rate":3.1,"Addr":"號, No. 1市府路信義區台北市台灣 110","Reviews":[{"StoreId":"ChIJP-nV3LmrQjQRfV4-XIaznRA","Text":"第一次因洽公來不知道要自備杯子，收銀服務員在解釋，此時一位女服務員說「你跟她說那麽多幹嘛？」.........妳們服務真的可以再好些吧？無言！！！！！"},{"StoreId":"ChIJP-nV3LmrQjQRfV4-XIaznRA","Text":"美式早餐好吃喔！"},{"StoreId":"ChIJP-nV3LmrQjQRfV4-XIaznRA","Text":"紅茶多樣又好喝，但市府內不提供紙杯，可用押金買或租杯子"},{"StoreId":"ChIJP-nV3LmrQjQRfV4-XIaznRA","Text":"10週年嘍! 咖啡好喝，太妃糖鮮奶茶也好喝"},{"StoreId":"ChIJP-nV3LmrQjQRfV4-XIaznRA","Text":"咖啡好喝"}]},{"Id":"ChIJh6rd0LWrQjQRQ2YnzBlJpys","Name":"Sunnyday Coffee \u0026 Select Shop/simple market","Rate":0,"Addr":"110台灣台北市信義區松勤街50號","Reviews":[]},{"Id":"ChIJzQAOKCwjaTQRWTUHYNd1-KM","Name":"Starbucks","Rate":3.9,"Addr":"No. 11, Song-Sho Road, Taipei City, 台北市台灣 105","Reviews":[{"StoreId":"ChIJzQAOKCwjaTQRWTUHYNd1-KM","Text":"坐位少 人潮多 店位小 但很方便"},{"StoreId":"ChIJzQAOKCwjaTQRWTUHYNd1-KM","Text":"坐位真的很少!\n人潮很多，建議不趕時間再去購買"},{"StoreId":"ChIJzQAOKCwjaTQRWTUHYNd1-KM","Text":"很小的分店，提供一個短暫休憩的空間。"},{"StoreId":"ChIJzQAOKCwjaTQRWTUHYNd1-KM","Text":"蠻小的 適合買了逛街"},{"StoreId":"ChIJzQAOKCwjaTQRWTUHYNd1-KM","Text":"No chairs to sit. Don't like that Starbucks. "}]},{"Id":"ChIJ_Y_2-7qrQjQRUCBvHVZMqkA","Name":"統一星巴克 松仁門市","Rate":3.7,"Addr":"110台灣台北市信義區松仁路7號","Reviews":[{"StoreId":"ChIJ_Y_2-7qrQjQRUCBvHVZMqkA","Text":"30 mins to get an Americano. Ridiculous. And bad service."},{"StoreId":"ChIJ_Y_2-7qrQjQRUCBvHVZMqkA","Text":"最近多了一個長桌，可以坐的位置又多了一些。門市裡沒有插座，大概是翻桌率考量。靠窗的位置能看到的景色也就是一般信義區。如果是一個人，不太建議，兩個人以上想聊天，不妨試試。"},{"StoreId":"ChIJ_Y_2-7qrQjQRUCBvHVZMqkA","Text":"A cosy Starbucks inside the Cathay United Bank building ground floor "},{"StoreId":"ChIJ_Y_2-7qrQjQRUCBvHVZMqkA","Text":"101可看台北，但全部為玻璃帷幕，夏天有點熱"},{"StoreId":"ChIJ_Y_2-7qrQjQRUCBvHVZMqkA","Text":"平日上班人潮多，座位蠻少的"}]},{"Id":"ChIJraeA2rarQjQRcsQVAszSNog","Name":"Starbucks 統一星巴克 (101 35F 門市)","Rate":3.9,"Addr":"110台灣台北市信義區信義路五段7號35樓之一","Reviews":[{"StoreId":"ChIJraeA2rarQjQRcsQVAszSNog","Text":"If you plan to go here be sure to make a reservation at least a day before. Ask the staff at the hotel you are staying or a person you know in Taiwan to make the reservation for you. We asked the host in AirBnB to make the reservation for us and she was really nice to do so.\n\nIt is nice and quite here. Very small though but it's fine. Also, come a little earlier or you'll lose your reservation. They also only allow a certain amount of people at a time.\n\nAlso, please don't wear slippers or shorts when you are inside their building nor should you go anywhere besides Starbucks.\n\nThe food and drinks aren't any different and they don't really have a certain special drink or snack which can be bought here. However, the good thing is you can get a good view of Taipei and don't have to pay a fee. Just the food you eat.\n\nThough I heard you can also get a good view from the Town hall for free. Haha."},{"StoreId":"ChIJraeA2rarQjQRcsQVAszSNog","Text":"服務態度不友善以外，環境很髒亂，尤其是地板上很多餅乾屑...讓人感受不到星巴克品牌價值的一個地方。"},{"StoreId":"ChIJraeA2rarQjQRcsQVAszSNog","Text":"週日正午，人潮不多。\n氣氛悠閒、服務友善，\n舒適地喝杯飲料吃個甜點，看看風景，一個愉快的週末。"},{"StoreId":"ChIJraeA2rarQjQRcsQVAszSNog","Text":"35층. 올라가기 위해서는 1층에서 스타벅스 직원을 기다려야 한다. 경치 좋은 자리를 위해 줄을 맨 앞에 서거나, 엘리베이터에 내려서 뛸 필요가 없다. 35층에 도착해서, 예약 번호 순서대로 입장(1층에서 포스트잇에 적어준 숫자가 예약 순서). 1인당 최소 200TWD 이상 구매해야됨. 조각케익이나 텀블러 사는 사람이 많지만, 텀블러가 필요없다면, stick (인스턴트) 커피 사면 됨. "},{"StoreId":"ChIJraeA2rarQjQRcsQVAszSNog","Text":"登高望遠😎 需要預約，每人有低消200元規定。服裝不能穿短褲和拖鞋。"}]},{"Id":"ChIJ2z8vxrCrQjQRQNOqzTQi4wM","Name":"Starbucks 統一星巴克 (信義ATT門市)","Rate":3.8,"Addr":"110台灣台北市信義區松壽路12號1樓","Reviews":[{"StoreId":"ChIJ2z8vxrCrQjQRQNOqzTQi4wM","Text":"早晨附近都還沒開始營業時，這家星巴克就先開始營業，很適合來喝杯咖啡"},{"StoreId":"ChIJ2z8vxrCrQjQRQNOqzTQi4wM","Text":"鬧中取靜，看晚電影或逛完百貨公司後的好去處。"},{"StoreId":"ChIJ2z8vxrCrQjQRQNOqzTQi4wM","Text":"服務人員親切，讓人一去再去"},{"StoreId":"ChIJ2z8vxrCrQjQRQNOqzTQi4wM","Text":"普通，整體上沒有太大的問題"},{"StoreId":"ChIJ2z8vxrCrQjQRQNOqzTQi4wM","Text":"店內較擁擠 但服務很好"}]},{"Id":"ChIJO8RPJ7erQjQRK7lBH_vJJQs","Name":"TWG Tea Salon \u0026 Boutique","Rate":4.2,"Addr":"110台灣台北市信義區市府路45號台北101購物中心 L5-01","Reviews":[{"StoreId":"ChIJO8RPJ7erQjQRK7lBH_vJJQs","Text":"Afternoon Tea Time -- a wonderful respite in Taipei 101. This place is absolutely beautiful. It is elegant and quiet and was the first time in a couple weeks of our Taiwan travels that we saw forks and knives! Tea menu is overwhelming! I enjoyed a superb Earl Grey and had this with their scones. Pricey but wonderful. We were relaxed and it was soothing."},{"StoreId":"ChIJO8RPJ7erQjQRK7lBH_vJJQs","Text":"驚訝這樣好的店，應該訓練有素的\n服務我們的漂亮小姐先是點錯餐，說她介紹錯誤，好沒關係就吃吧。\n再上錯茶，倒錯別人的茶給我們，小姐再將鼻子湊近茶壺口聞確定有錯，好吧只能說可能下午較忙。\n現場其他的服務確實不錯。"},{"StoreId":"ChIJO8RPJ7erQjQRK7lBH_vJJQs","Text":"Excellent staff who really go out their way to assist you. Food was delicious and we would highly recommend visiting."},{"StoreId":"ChIJO8RPJ7erQjQRK7lBH_vJJQs","Text":"好茶配好點心 \n在101內的金色裝潢裡顯現茶的獨特\n整ㄨ的茶罐跟印著滿滿茶名的點茶單\n哇！ 需要有人來好好介紹一下囉⋯⋯"},{"StoreId":"ChIJO8RPJ7erQjQRK7lBH_vJJQs","Text":"下午茶的馬芬很好吃，茶類很多種但服務人員會幫忙推薦，用餐時間沒有限制可以安心的享受悠閒的時光。"}]},{"Id":"ChIJjccWxrarQjQRm90LTHdMHnQ","Name":"Cafe Lugo","Rate":3.4,"Addr":"號 B1( 101), No. 45市府路信義區台北市台灣 110","Reviews":[{"StoreId":"ChIJjccWxrarQjQRm90LTHdMHnQ","Text":"Nice cafe with quite friendly service. The waffles were not spectacular, it was the rubbery eggy type and not the fluffy type. The cafe latte was nice. The deco is unique with coffee pots hanging from the ceiling."},{"StoreId":"ChIJjccWxrarQjQRm90LTHdMHnQ","Text":"付完錢沒找錢服務人員就飄走了\n然後咖啡等了四五十分鍾都沒好\n最後知道他們疑似漏單了 就退錢走人\n還能退錢還不錯\n員工的教育訓練有待加強"},{"StoreId":"ChIJjccWxrarQjQRm90LTHdMHnQ","Text":"令人驚訝拿鐵比星巴克濃醇，不過沒有提供低脂肪牛奶可選。。。。冰滴也不夠濃厚，炭燒味偏重，有插座可使用電腦。"},{"StoreId":"ChIJjccWxrarQjQRm90LTHdMHnQ","Text":"餐點尚可，咖啡ok，是逛101逛累了可以來喝杯咖啡休憩一下的地方~"},{"StoreId":"ChIJjccWxrarQjQRm90LTHdMHnQ","Text":"平日下午來還蠻安靜的，可以咖啡和一個人靜空看書，五六日人多較吵雜，就比較不推薦了"}]},{"Id":"ChIJt0_SB7GrQjQRnVwE2yjwnXk","Name":"Cama Café - 台北松仁店","Rate":3.9,"Addr":"110台灣台北市信義區松仁路32-36號","Reviews":[{"StoreId":"ChIJt0_SB7GrQjQRnVwE2yjwnXk","Text":"在信義計畫區的小小角落擠了這間小小的店，提供這裡的上班族平價質優的咖啡，其實是很感恩的…"},{"StoreId":"ChIJt0_SB7GrQjQRnVwE2yjwnXk","Text":"小小一家店，很可愛。"},{"StoreId":"ChIJt0_SB7GrQjQRnVwE2yjwnXk","Text":"咖啡一般般，店員似乎不想做生意，還直接在店外頭抽菸，希望他們再加油。"},{"StoreId":"ChIJt0_SB7GrQjQRnVwE2yjwnXk","Text":"\u003e 小小一間藏在辦公大樓的門市\n\u003e 咖啡品質就如一般Cama的水準\n\u003e 空調似乎常常沒開，有點悶熱"},{"StoreId":"ChIJt0_SB7GrQjQRnVwE2yjwnXk","Text":"附近買咖啡的好地方"}]},{"Id":"ChIJh9DVybCrQjQR9sf9bMkkgKA","Name":"Isaac Toast \u0026 coffee","Rate":2.1,"Addr":"110台灣台北市信義區松壽路12號4F","Reviews":[{"StoreId":"ChIJh9DVybCrQjQR9sf9bMkkgKA","Text":"狠雷的味道 , 完全跟韓國吃到的味道完全不同 , 店員動作超慢 , 我點了餐等了10分鐘才拿到 , 二個人一起泡咖啡 , 二個人一起整理東西 \n只有我一個客人 , 好失敗的土司"},{"StoreId":"ChIJh9DVybCrQjQR9sf9bMkkgKA","Text":"吐司太濕了 和韓國的店差得有點多\n動作稍顯緩慢 但應該是動線流動排的不順暢啦\n如果要吃還是建議去韓國吃😝"},{"StoreId":"ChIJh9DVybCrQjQR9sf9bMkkgKA","Text":"明明就沒人，還慢到不行，兩個店員邊做邊聊天就算了，包三明治時客人就在面前還在慢慢包慢慢聊，無敵傻眼。工讀生的素質可以挑一下嗎?\n完全不專業"},{"StoreId":"ChIJh9DVybCrQjQR9sf9bMkkgKA","Text":"吐司本身還不錯，但裡面的料一場悲劇……特別是豬肉……\n當年從韓國回來後才紅，沒吃過，特地前來朝聖\n結果……美而美還比較好吃 至少它便宜😂\n"},{"StoreId":"ChIJh9DVybCrQjQR9sf9bMkkgKA","Text":"松山店今天開幕，前去朝聖，取餐經驗不太愉快，盼望能改進\n因為要等很久就先去別處逛之後回來一定過號\n採現場叫號所以沒過號的幾乎都不會排隊\n但在排隊路線卻看到第一個人在跟店員聊天，我知道該店員很為難因為很忙，但沒阻止客人，之後排隊的第二位客人看對方沒要取餐就過去領，第三位是我和我後面的也都是過號在排隊...\n結果輪到我的時候那位聊天的客人突然擋在我前面，因為她的號碼才到，既然他排在前頭讓他領也沒有什麼\n就在這時，突然有人從後面擠過來，說: 我要趕火車 到幾號了?\n店員說才正要起鍋，請他等一等，然後全都在忙包裝，於是乎我們這幾個排隊的又在傻等\n但我想給別人方便就是給自己方便，就等吧\n這時候我後面客人的朋友回來了，很大聲的說怎麼還沒好? 只看到店員就邊忙稍微看了一下他們，然後又繼續忙自己的事\n我後面的客人只好說好像號碼還沒到吧...\n接著趕火車的人終於拿到東西走了，就看店員低頭一直忙東西(動作很慢)\n我想說我應該要問了，但他還是很忙自顧自的在包就想說觀察一下他在做什麼\n包好了之後他開始叫號，結果都沒人理他\n然後他又想要繼續低頭包東西\n這時候我恍然大悟...原來排隊是排假的，一定要直接跟他說才行\n而且叫號後我後面的客人也才知道他過號了\n如果不是由我來叫他 那排隊就真的是排假的\n所以我終於叫他了，我說為什麼這裡在排隊，還是一堆人從後面擠過來領東西，你們都完全沒在管我們排隊的\n他還跟我說我們這是叫號啊\n我不想再說什麼 只有說: 你叫了現在沒有人要來拿 那我可以取餐了嗎?\n店員看了我的號碼單並且神速的把我的東西交給我說: 不好意思\n我只有回聲: 恩  (我也覺得不太禮貌 但我也是第一次被激到)\n\n本來不想反應的 因為大概不會想再去\n如果是第一天開幕比較亂 盼望能改進"}]},{"Id":"ChIJrSZypLCrQjQRzeaix7XEgA8","Name":"Nescafe Dolce Gusto 信義新光三越A9專櫃","Rate":0,"Addr":"110台灣台北市信义区松壽路9號4樓","Reviews":[]},{"Id":"ChIJmbOzNbqrQjQR1Q5INUqTDHg","Name":"Nespresso 新光三越 台北信義新天地A8館 專櫃","Rate":3,"Addr":"110台灣台北市信義區松高路12號7樓","Reviews":[{"StoreId":"ChIJmbOzNbqrQjQR1Q5INUqTDHg","Text":"櫃姐漂亮服務又親切😍"},{"StoreId":"ChIJmbOzNbqrQjQR1Q5INUqTDHg","Text":"杯子很髒"}]},{"Id":"ChIJsxIIdrqrQjQRWAJUjw9faBc","Name":"Nespresso 新光三越 台北信義新天地A11館 精品店","Rate":4.4,"Addr":"110台灣台北市信義區松壽路11號","Reviews":[{"StoreId":"ChIJsxIIdrqrQjQRWAJUjw9faBc","Text":"整間店充滿了時尚感、 最重要的是售後服務非常的用心、給矛滿分鼓勵！"},{"StoreId":"ChIJsxIIdrqrQjQRWAJUjw9faBc","Text":"產品、風格、服務、價格全部都很到位\n尤其在價格調整後單顆16元的均價非常值得"},{"StoreId":"ChIJsxIIdrqrQjQRWAJUjw9faBc","Text":"來自巴西的咖啡，膠囊機器和咖啡膠囊及周邊商品，24種不同濃郁的口感。可以排隊試飲，服務非常專業，介紹很仔細哦！鋁製膠囊可以回收，咖啡渣也可以再次利用，真的不錯！為自己一分鐘泡杯咖啡吧！"},{"StoreId":"ChIJsxIIdrqrQjQRWAJUjw9faBc","Text":"排隊就可以免費試喝各種咖啡，服務人員都願意給予各種建議\n因為這個店是唯一可以看到咖啡膠囊機器及膠囊販售的店\n要買限定款，不用網路通路，到此就可以試喝及購買\n雖然屬於百貨公司的一角，但是前後都有門，方便出入"},{"StoreId":"ChIJsxIIdrqrQjQRWAJUjw9faBc","Text":"服務人員專業，等待時間可品嚐多種口味咖啡。"}]},{"Id":"ChIJfxsYC7qrQjQRl_C2wgvWAj0","Name":"喬尼亞咖啡","Rate":0,"Addr":"110台灣台北市信義區松壽路9號(信義連通空橋)4樓","Reviews":[]},{"Id":"ChIJcxkKdLqrQjQRfnJrOy6VUkc","Name":"AMP Café 新光三越 A11","Rate":5,"Addr":"No. 11, Songshou Road B2, 信義區台北市台灣 110","Reviews":[{"StoreId":"ChIJcxkKdLqrQjQRfnJrOy6VUkc","Text":"咖啡很棒, 外帶方便, 店員親切, 值得再去！"}]},{"Id":"ChIJt6RubLqrQjQRa9ARszu1CNk","Name":"LINE FRIENDS Cafe \u0026 Store","Rate":4.2,"Addr":"110台灣台北市信義區松壽路11號新天地A11館","Reviews":[{"StoreId":"ChIJt6RubLqrQjQRa9ARszu1CNk","Text":"超級可愛❣️東西也很好吃呢😋"},{"StoreId":"ChIJt6RubLqrQjQRa9ARszu1CNk","Text":"只有可愛可以形容～隔壁有主題咖啡廳，只能外帶，但旁邊有立食區，主要賣飲料跟雞蛋糕。"},{"StoreId":"ChIJt6RubLqrQjQRa9ARszu1CNk","Text":"很好玩有趣，幾乎所有line家族周邊商品都有賣的地方，\n店面佔地超大，很多實體line娃娃雕像可以一起拍照。"},{"StoreId":"ChIJt6RubLqrQjQRa9ARszu1CNk","Text":"熊大妹好可愛，可惜有些東西臺灣沒有賣，像是觸碰燈和和巴黎熊大兔兔手機套"},{"StoreId":"ChIJt6RubLqrQjQRa9ARszu1CNk","Text":"很可愛的地方，有一隻很大的熊大可以拍照，也有VR實境可以體驗，就在新光三越A11，逛完百貨公司不妨來走走吧"}]},{"Id":"ChIJrSZypLCrQjQR2NYLSstifcw","Name":"Caffé Florian 福里安花神咖啡館","Rate":3.1,"Addr":"110台灣台北市信義區松壽路9號新光三越信義新天地A92F","Reviews":[{"StoreId":"ChIJrSZypLCrQjQR2NYLSstifcw","Text":"[如果你想花大錢當傻子 這間餐廳絕對符合你想被坑殺的期待]\n因為威尼斯總店的美好回憶 讓我們決定踏入位在台灣的分店 事實上這個決定就是災難的開始:\n\n(1) 喧鬧擁擠環境：\n光是我們點的早餐盤就佔滿座位 桌與桌的空間也非常狹窄 導致隔壁桌聊天聲字字句句猶言在耳 環境吵雜不堪 (當然台灣人民的素質也要檢討 不管再高級的餐廳總是喜歡大聲嚷嚷 讓人惱火)\n\n(2) 菜單圖片不實：\n實際份量跟菜單上的示意圖「差異極大」難怪我們點含服務費近台幣一千的早餐盤 「只符合一人低消」必須再加點一杯飲料  感覺就只是吃高級餐具的氣氛 而非食物本身 食物份量真的少的非常誇張\n\n(3) 義式日常輕食：\n早餐盤內容物的確是義大利不錯等級的火腿/可頌/起司 但是 這在當地3星級以上的飯店是很基本的早餐規格 味道完全一模一樣 還以為名店味道更好\n\n(4) 現做時間過長：\n一個偏「冷盤」的早餐套餐 告知現做要30-40分鐘 令人不敢置信 我們在義大利也沒如此久候\n\n(5) 服務品質差勁：\n等候的漫長歲月中 竟然連一杯水都沒端上 服務員就完全消失地無影無縱 好不容易找到後 告知如此誇張的服務還被臭臉以對 待水送上來時 我們忍住內心的不悅 保持應有的禮貌向其道謝 服務人員此時竟然靜默 彷彿是怨嘆我們找他麻煩 令人匪夷所思(後來換成不同服務人員 接待明顯改善 只是仍舊一張撲克臉)\n\n總結：以這間高定價與高名氣的餐廳 如此的服務/食物品質非常不及格 義大利人平時那套的孤傲自我雖然令人退避三舍 但至少提供合理的餐飲品質 無奈輸出到台灣 好的不學 壞的到是學得唯妙唯肖！？ (事實上威尼斯總店的服務人員不僅接待客人時面帶微笑 還彬彬有禮)"},{"StoreId":"ChIJrSZypLCrQjQR2NYLSstifcw","Text":"很舒服的一間店,服務好氣氛佳，有次一進去覺得太吵馬上就離開了。 服務生送水送的非常勤，常常一壺都還沒喝完就送上下一壺。經典巧克力,皇家咖啡和提拉米酥很棒，沒吃過那麼好吃的點心！若沒人在吵的話，是間值得拜訪的店"},{"StoreId":"ChIJrSZypLCrQjQR2NYLSstifcw","Text":"來自義大利的福里安花神咖啡館，終於有機會來朝聖啦～\n店鋪開設在新光三越的A9館2樓，從市政府捷運站走過來著實需要花點時間，還好本姑娘腳程快，一個箭步搶先來櫃檯登記呀！\n\n很幸運的30分鐘內就接到電話啦，開心去報到的同時也看到門口已經開始排隊啦（O.S.不用排隊就是～爽～）\n\n店內低消是飲品or冰淇淋，我們捨棄了下午茶組合，採單點「義式馬鈴薯鹹派、挪威鮭魚三明治、英式早餐茶、福里安290週年紀念版咖啡及台灣限定版的茉莉花茶慕斯蛋糕」\n\n餐點部分，鹹派讓人頗為失望（吃下肚只覺得就是普通的馬鈴薯），三明治的口感很好（軟軟的吐司搭配鮭魚讓人想一口接著一口），紀念版咖啡（比想像中的還要小杯～另外因為有添加酒～所以怕酒味的人不建議點選），茶品部分，並非所有菜單上面都茶都有（因為都是從義大利空運過來，所以賣完就要再等下次進貨啦）至於最後的茉莉花茶慕斯（❤️台灣限定、上頭是帶有茉莉花茶香氣的慕斯～入口後會有淡淡的茉莉花香～下層配上鮮奶油蛋捲～奶油搭上海綿蛋糕口感極佳，是款極為推薦的甜品❤️）\n\n整體而言，店鋪空間還不錯，幸運的被安排到沙發區（沙發軟軟的很好座），餐廳內設有廁所（空間很大、不用特地離開餐廳去。\n\n消費需要加上一成！"},{"StoreId":"ChIJrSZypLCrQjQR2NYLSstifcw","Text":"如果你錢太多，剛好在這想找地方休息，可以來這家。\n\n東西很貴，咖啡不好喝，蛋糕不錯吃，但我寧願去我家附近蛋糕店買蛋糕，還不到100。\n\n我，不會，再，來。"},{"StoreId":"ChIJrSZypLCrQjQR2NYLSstifcw","Text":"全球10大最美咖啡館 - Caffe Florian\n\n在炙熱的午后與\n威尼斯人最愛...優雅浪漫下午茶 相遇\n代表優雅的 華麗銀製大托盤上場(義大利原裝)...\n\n還是我的最愛 #水牛起司可頌\nFlorian俱樂部三明治\nCoppa Caffe Florian -\n咖啡 + 提拉米蘇 + 巧克力冰淇淋 + 香甜酒 + 巧克力醬 +鮮奶油\n\n聖代... 因為有酒.. 所以我愛"}]},{"Id":"ChIJVXn7qLurQjQRCDilM630zr8","Name":"Mauink 墨癮精品咖啡豆專賣","Rate":0,"Addr":"號 B2, No. 28松仁路信義區台北市台灣 110","Reviews":[]},{"Id":"ChIJt1_LrrCrQjQRoe7Z1sIjVjk","Name":"Chôn Select Store \u0026 Cafe","Rate":5,"Addr":"110台灣台北市信義區松壽路3-1號號","Reviews":[{"StoreId":"ChIJt1_LrrCrQjQRoe7Z1sIjVjk","Text":"Fantastic food and amazing coffee! Staff really knows there products and are happy to serve! "},{"StoreId":"ChIJt1_LrrCrQjQRoe7Z1sIjVjk","Text":"非常專業的手沖咖啡店，還可受到專業的裁判級專家的指導。"}]},{"Id":"ChIJbZ_CVaSrQjQRvwHr-syXe24","Name":"嗨咖咖啡手沖咖啡館","Rate":4.8,"Addr":"110台灣台北市信義區虎林街164巷78號","Reviews":[{"StoreId":"ChIJbZ_CVaSrQjQRvwHr-syXe24","Text":"相當讓人放鬆的咖啡廳，地點因為偏僻所以相當安靜，鄰近一高級住宅與公園，因此景致還算不錯。單品咖啡部分有很多種地區的可以選擇，調味咖啡部分選擇也不少，就算不點咖啡也還有奶茶與藥草茶可以選擇。甜點部分選擇略少，不過就個人嚐到的蛋糕而言相當合胃口，可以期待其他的甜點應該都不錯。\n另外店內附有一書架，可以自由在店內閱讀上面的書籍，就算不想看書，店內也有 Wi-Fi 可供使用。"},{"StoreId":"ChIJbZ_CVaSrQjQRvwHr-syXe24","Text":"從流浪類型的單車咖啡開始追，這禮拜總算開了店面！\n\n嗨咖咖啡手沖咖啡一直是我的最愛，老闆總是很樂意跟客人們分享各咖啡豆的精妙之處，不強迫推銷，就只為了找出客人最喜歡的精品滋味。\n\n店內還新增了以前沒有的義式濃縮，拿鐵、卡布、焦糖拿鐵、奶茶等飲料品項，每一杯都很優秀！\n\n店面在虎林街一個小公園旁，店外廣場林蔭遮蔽，坐在店內外，來杯咖啡都是莫大的享受！\n\n非常推薦！\n\nOne of the best poured over coffee in this area! Cozy place with good coffee, vest la vie !\n\n"},{"StoreId":"ChIJbZ_CVaSrQjQRvwHr-syXe24","Text":"A lovely and cozy coffee shop next to my building. Only one barista here. The latte is really nice. "},{"StoreId":"ChIJbZ_CVaSrQjQRvwHr-syXe24","Text":"This specialty coffee spot is amazing. We had iced pour-overs and sampled the brown sugar pastry that the owner made in the back. Perfect. "},{"StoreId":"ChIJbZ_CVaSrQjQRvwHr-syXe24","Text":"環境舒適，老闆對咖啡很專業很有堅持，動物友善。"}]},{"Id":"ChIJSTLZ6barQjQR7OeldE1QDsU","Name":"LALOS Bakery (101店)","Rate":4.1,"Addr":"110台灣台北市信義區市府路45號B1樓","Reviews":[{"StoreId":"ChIJSTLZ6barQjQR7OeldE1QDsU","Text":"店員服務的態度很差\n我用apple pay 剛開始沒有把指紋放上去經過她提點放上去付款完成，她也不知道什麼意思推了我的手機，或許是要我把卡片移開，但我沒見過這麼沒禮貌的動作。結完帳也沒有標準的話術：謝謝\n默默地把麵包跟收據就推給我，完成交易\n看來他們是不缺客人"},{"StoreId":"ChIJSTLZ6barQjQR7OeldE1QDsU","Text":"試吃蠻大方，有變化的拖鞋麵包有小驚艷到，價錢在101算便宜了。服務的確比較冷淡一些。"},{"StoreId":"ChIJSTLZ6barQjQR7OeldE1QDsU","Text":"麵包極有嚼勁，蕎麥長棍、經典裸麥圓麵包都沒有過多調味但十分美味。深得歐式麵包精髓。"},{"StoreId":"ChIJSTLZ6barQjQR7OeldE1QDsU","Text":"怎麼會！太厲害的麵包了，單純的美好！吃後令人忘卻不了的口感，唇齒留香😍怎麼會有那麼好吃的麵包！！\n\n可惜是結帳速度很日本式，要稍微等待，但是得值得！"},{"StoreId":"ChIJSTLZ6barQjQR7OeldE1QDsU","Text":"麵包用料實在，講究法國進口，大廚也是法國人！ 實在非常喜歡"}]},{"Id":"ChIJUztrxbWrQjQRGsdUosiD8LQ","Name":"好丘 Good Cho's","Rate":4,"Addr":"110台灣台北市信義區松勤街54號","Reviews":[{"StoreId":"ChIJUztrxbWrQjQRGsdUosiD8LQ","Text":"貝果好吃😋"},{"StoreId":"ChIJUztrxbWrQjQRGsdUosiD8LQ","Text":"東西很貴"},{"StoreId":"ChIJUztrxbWrQjQRGsdUosiD8LQ","Text":"平日不用等待，但服務員品質不一，別桌介紹得很詳盡，我們第一次來但是很隨便叫我們自己看，然後連水都吃到一半才補上。\n\n如果能被專業服務員服務到，會是不錯的用餐體驗，希望能加強訓練，讓每個客人都能享有相同的服務品質。"},{"StoreId":"ChIJUztrxbWrQjQRGsdUosiD8LQ","Text":"Nice cosy place with healthy but tasty food. Very spacious and relaxing atmosphere.   It's a pity that the air conditioning is not strong enough for our hot weather. "},{"StoreId":"ChIJUztrxbWrQjQRGsdUosiD8LQ","Text":"服務品質有待加強。現場留了電話候位通知，但是後面再去問的位子的人，也沒訂位直接放進去是什麼意思？"}]},{"Id":"ChIJQcKOdbqrQjQROgrX00ZQqmo","Name":"新光三越台北信義新天地A11","Rate":4.2,"Addr":"110台灣台北市信義區松壽路11號","Reviews":[{"StoreId":"ChIJQcKOdbqrQjQROgrX00ZQqmo","Text":"Very cold air con and good department store. Fantastic design of building and interior. Restaurant is also my favorite to go."},{"StoreId":"ChIJQcKOdbqrQjQROgrX00ZQqmo","Text":"動線標示不夠明顯，建議棟別標示要明確清楚一些。"},{"StoreId":"ChIJQcKOdbqrQjQROgrX00ZQqmo","Text":"不錯，挺好的，可以逛逛"},{"StoreId":"ChIJQcKOdbqrQjQROgrX00ZQqmo","Text":"悠閒吹冷氣的逛街很舒服"},{"StoreId":"ChIJQcKOdbqrQjQROgrX00ZQqmo","Text":"眼花撩亂，淡定散心"}]},{"Id":"ChIJG2UmvbCrQjQRNaSZwjk2gZ0","Name":"Krispy Kreme","Rate":3.7,"Addr":"110台灣台北市信義區松壽路18號","Reviews":[{"StoreId":"ChIJG2UmvbCrQjQRNaSZwjk2gZ0","Text":"標準美式甜甜圈，甜度很夠，口味繁多，但推薦原味。"},{"StoreId":"ChIJG2UmvbCrQjQRNaSZwjk2gZ0","Text":"原味甜甜圈最好吃❤"},{"StoreId":"ChIJG2UmvbCrQjQRNaSZwjk2gZ0","Text":"不難吃但是有一點甜～喜歡巧克力口味跟原味\n位於熱鬧區域還蠻好找的"},{"StoreId":"ChIJG2UmvbCrQjQRNaSZwjk2gZ0","Text":"If you aspire to grow into an American + size, this is the right place to do it. Each donut is carefully designed to deliver more than 30 times the sugar your body needs daily. But man, it's good."},{"StoreId":"ChIJG2UmvbCrQjQRNaSZwjk2gZ0","Text":"原味超級無敵好吃，外酥內軟加上外層薄薄的糖霜，只能用美味來形容了！"}]},{"Id":"ChIJjfzEsLCrQjQRsd6k1mO2-WA","Name":"ATT 4 FUN 臺北信義店","Rate":4.1,"Addr":"110台灣台北市信義區松壽路12號","Reviews":[{"StoreId":"ChIJjfzEsLCrQjQRsd6k1mO2-WA","Text":"喜歡4樓甜點區。\n離世貿中心近，做展覽和看展覽的人可以去走走覓食。"},{"StoreId":"ChIJjfzEsLCrQjQRsd6k1mO2-WA","Text":"甜點店，韓式料理"},{"StoreId":"ChIJjfzEsLCrQjQRsd6k1mO2-WA","Text":"購物新歡地點！"},{"StoreId":"ChIJjfzEsLCrQjQRsd6k1mO2-WA","Text":"Ok la"},{"StoreId":"ChIJjfzEsLCrQjQRsd6k1mO2-WA","Text":"😎😎"}]},{"Id":"ChIJUyWwC7qrQjQRj7g1LFj7SFg","Name":"新光三越台北信義新天地A9","Rate":4.2,"Addr":"110台灣台北市信義區松壽路9號","Reviews":[{"StoreId":"ChIJUyWwC7qrQjQRj7g1LFj7SFg","Text":"附近百貨商場餐廳與影城林立,逛完需要時間,停車場停松壽公園地下停車場較便宜,版或公司停車費貴但消費可抵停車費"},{"StoreId":"ChIJUyWwC7qrQjQRj7g1LFj7SFg","Text":"商店，餐廳處處都有，空橋可以直接通往四個館，不必走到一樓，難怪人潮多"},{"StoreId":"ChIJUyWwC7qrQjQRj7g1LFj7SFg","Text":"現場的櫃檯不曉得在跩什麼，好聲好氣的跟妳說忘記帶收據，第一次來消費誰天生知道妳們的規矩啊，而且結帳的時候也沒有提醒要帶單子，客氣的跟妳說忘記帶要怎麼辦，還瞪我，這是怎麽一回事啊？"},{"StoreId":"ChIJUyWwC7qrQjQRj7g1LFj7SFg","Text":"非常小资的一个地方，就在台北101附近，灯光有淮海路的感觉。优质高端的购物中心。"},{"StoreId":"ChIJUyWwC7qrQjQRj7g1LFj7SFg","Text":"還算方便的百貨公司，有我要的櫃。汽車也好停。\n挑人少時來買 才不擠。"}]},{"Id":"ChIJ5zCG5a-rQjQRSbC2PiZOWKY","Name":"WOW FURNITURE 北歐進口家具","Rate":4.2,"Addr":"110台灣台北市信義區信義路五段91巷26號","Reviews":[{"StoreId":"ChIJ5zCG5a-rQjQRSbC2PiZOWKY","Text":"精緻的家具店"},{"StoreId":"ChIJ5zCG5a-rQjQRSbC2PiZOWKY","Text":"家居"},{"StoreId":"ChIJ5zCG5a-rQjQRSbC2PiZOWKY","Text":""},{"StoreId":"ChIJ5zCG5a-rQjQRSbC2PiZOWKY","Text":""}]},{"Id":"ChIJO7i2NrqrQjQRw_HGC_L9VKY","Name":"Kiehl's - 新光三越A8","Rate":0,"Addr":"110台灣台北市松高路12號","Reviews":[]},{"Id":"ChIJIXvYsLCrQjQRWsrfQUG1iGg","Name":"ZARA HOME","Rate":4.1,"Addr":"11051台灣台北市信義區松壽路12號","Reviews":[{"StoreId":"ChIJIXvYsLCrQjQRWsrfQUG1iGg","Text":"生活不能沒有Zara Home，推薦床單床包等織品以及擴香類商品，超棒！！"},{"StoreId":"ChIJIXvYsLCrQjQRWsrfQUG1iGg","Text":"東西品質挺穩定的，簡潔又漂亮，裡面也放著香芬讓人逛起來心情很好\""},{"StoreId":"ChIJIXvYsLCrQjQRWsrfQUG1iGg","Text":"不同於其他ZARA ，這家販賣的是家居用品如傢具、床、檯燈…"},{"StoreId":"ChIJIXvYsLCrQjQRWsrfQUG1iGg","Text":"每次經過都會被香氛吸引~~~販售的傢飾品都很有質感，只是也不便宜。"},{"StoreId":"ChIJIXvYsLCrQjQRWsrfQUG1iGg","Text":"家具不算差，跟無印量品走不同的風格。"}]},{"Id":"ChIJxQqtNbqrQjQR3iHT3CHwcWA","Name":"Georg Jensen","Rate":0,"Addr":"松高路12號1F 信義區, 台北市, 台灣 110","Reviews":[]},{"Id":"ChIJ7RVlNLqrQjQRHx6mfwOkBEU","Name":"Georg Jensen (at 101 Mall)","Rate":0,"Addr":"市府路45號1樓 信義區, Sinyi District, 台北市, 台灣 110","Reviews":[]},{"Id":"ChIJiXSSILerQjQRuq_D9uHVFMg","Name":"POLO RALPH LAUREN (101店)","Rate":4.5,"Addr":"110台灣台北市信義區市府路45號","Reviews":[{"StoreId":"ChIJiXSSILerQjQRuq_D9uHVFMg","Text":"服務很好，樣式很多😍😘😗"},{"StoreId":"ChIJiXSSILerQjQRuq_D9uHVFMg","Text":""}]},{"Id":"ChIJ936ujy7RDRQRh3NS--5UCgM","Name":"微風廣場松高店","Rate":3.9,"Addr":"110台灣台北市信義區松高路16號","Reviews":[{"StoreId":"ChIJ936ujy7RDRQRh3NS--5UCgM","Text":"這裡真是太好買!snoopy服飾因為是made in TAIWAN所以比在日本買便宜,遇到店家有活動特惠時一模一樣的馬克杯竟比在日本買還要便直!!"},{"StoreId":"ChIJ936ujy7RDRQRh3NS--5UCgM","Text":"小小長長的，搶市佔率高。"},{"StoreId":"ChIJ936ujy7RDRQRh3NS--5UCgM","Text":"附近百貨林立，商品無特色"},{"StoreId":"ChIJ936ujy7RDRQRh3NS--5UCgM","Text":"好地方，值得來逛逛"},{"StoreId":"ChIJ936ujy7RDRQRh3NS--5UCgM","Text":"👍👍👍"}]},{"Id":"ChIJERAktLCrQjQRnmfBSUiorZQ","Name":"OASIS","Rate":0,"Addr":"110台灣台北市信義區松壽路12號","Reviews":[]},{"Id":"ChIJfY_HnrCrQjQRmIoMFKYGW7g","Name":"GLAM AIR","Rate":4.4,"Addr":"11051台灣台北市信義區松壽11號B1","Reviews":[{"StoreId":"ChIJfY_HnrCrQjQRmIoMFKYGW7g","Text":"Ingenious combination of cotton candy and high quality ice cream. Love! Quick and friendly service 👍🏼. Best order to-go as the set up is of a kiosk, with only 4 high stools at a tiny counter. "},{"StoreId":"ChIJfY_HnrCrQjQRmIoMFKYGW7g","Text":"推薦「璀璨葡萄柚」、「彩虹棉花糖霜淇淋」、Bling Bling銀河系列星空飲料【土星】-玫瑰葡萄口味及I Can Fly雲朵霜淇淋-檸檬蘇打口味很潮喔!"},{"StoreId":"ChIJfY_HnrCrQjQRmIoMFKYGW7g","Text":"棉花糖配冰淇淋，超好吃的搭配\n但跟台南餓魚咬冰比起來價格跟餐點差很多\n冰跟棉花糖都減半價錢也多一倍\n可是變化比較多有彩虹照型等等\n棉花糖也有加跳跳糖吃起來很驚奇\n推薦吃烏雲是薄荷口味吃起來比較不會甜膩\n它們還有飲品星空氣水很有特色加了亮亮的東西\n看起就像銀河很美"},{"StoreId":"ChIJfY_HnrCrQjQRmIoMFKYGW7g","Text":"來訪兩次，每次都讓我有不同驚喜!!我覺得GLAMAIR非常具創新能力，將棉花糖和冰品結合，更會做成許多造型，成功讓饕客吸睛，來這裡吃冰又可以逛街，一舉兩得，假日可以來這裡吹吹冷氣吃吃冰，分享給大家。"},{"StoreId":"ChIJfY_HnrCrQjQRmIoMFKYGW7g","Text":"I was pleasantly surprised! Thought this place would be all gimmick, no substance all bark and no bite but I was wrong. The ice cream was very creamy and silky and melted well with the cotton candy. They also put pop rocks in the ice cream... how amazing is that. Gonna put that in all my ice cream from now on. Cool spot. Would come back when in need of a sugar high."}]},{"Id":"ChIJbcUIeLGrQjQR7BC8o2k-LpQ","Name":"Queen \u0026 Daddy 珠寶會館","Rate":0,"Addr":"110台灣台北市信義區松仁路147號","Reviews":[]},{"Id":"ChIJScbdi7qrQjQRS-GbgQ4BkNc","Name":"Young Living Taiwan 悠樂芳","Rate":4.5,"Addr":"8th Floor, No. 89松仁路信義區台北市台灣 110","Reviews":[{"StoreId":"ChIJScbdi7qrQjQRS-GbgQ4BkNc","Text":"Probably the best will call I have seen for Young Living. It was very green and relaxing. Don't forget to go out to the personal balcony to get an amazing shot of Taipei 101."},{"StoreId":"ChIJScbdi7qrQjQRS-GbgQ4BkNc","Text":""},{"StoreId":"ChIJScbdi7qrQjQRS-GbgQ4BkNc","Text":""}]},{"Id":"ChIJ1TMC6barQjQRy0UmEQwFcvg","Name":"Apple Store 台北101","Rate":3.7,"Addr":"市府路45號, 臺北101, 信義區, 台北市, 台灣 110","Reviews":[{"StoreId":"ChIJ1TMC6barQjQRy0UmEQwFcvg","Text":"維修經驗非常不好，客服與現場人員說法完全不同，現場人員口氣非常不好的教訓我們，但明明已經先打客服詢問過，且文件聯絡完全不能使用email，不敢相信我買了這麼多科技產品還需要本人親送文件多次，無言至極。"},{"StoreId":"ChIJ1TMC6barQjQRy0UmEQwFcvg","Text":"蘋果迷必來的景點，門外有大大的蘋果logo，店內空間非常的大，燈光明亮又不刺眼，服務人員態度都非常熱情，很值得來這裡走走。"},{"StoreId":"ChIJ1TMC6barQjQRy0UmEQwFcvg","Text":"太棒了👏👏👏"},{"StoreId":"ChIJ1TMC6barQjQRy0UmEQwFcvg","Text":"一共打了genius bar四次。第一次和第四次的處理人員能專業處理並告知您該預先做的動作。某部分電話的技術支援不敢恭維，而且常等待（聽音樂）轉接很久！\n（若是有空，要有多餘時間等候處理，因爲上一位顧客可能狀況比較多，或許自己也會想多瞭解手機狀況而delay下一個顧客的維修時間。）推薦大家預約genius bar到101直營店，有sop流程還是會儘速處理！（但預約最好提早一週較有可能有自己可以的時間）\n雖然得親自且專程跑來（最好自己先做備份，因為現場多人共用wifi，所以速度有限。），但現場人員面對面處理，能有效率處理。幸運遇到sunny,電洽遇上Anita,最後也遇上另一個專業的技術人員，維修完畢還貼心幫我收了耳機線。\n真的感謝在電話支援求救極度絕望之際，遇上用心專業的genius bar小姐幫我轉接洽101很細心的Anita，Anita很有耐心的回撥了很多次電話與我聯繫並確定二次維修!\n總之，第一次使用iphone並且維修iphone的過程雖不簡單，但暫時落幕！看著技術人員秉著對iphone的熱愛，下次會接著使用iphone 的！：）"},{"StoreId":"ChIJ1TMC6barQjQRy0UmEQwFcvg","Text":"服務系統有問題，只能採預約制，如果是為了保証有人服務的品質那可以理解，但在預約的時間報到後，還是繼續等，過了快一個小時還是沒有人來處理，問工作人員也沒辦法預估還要等多久，只能繼續等待...續\n\n筆電檢修結果是電池壞了，整個過程還得跑三趟，趟趟要預約：\n1)檢查確認原因，要申請維修零件(電池)，需1~2個禮拜，產品請自行帶回！\n2)八天後通知零件已備妥，需在期限前(五天，含星期六日)將產品再帶去．隔兩天帶去，筆電留下，簽了切結書，然後被告知還要3~5天，等候通知．\n3)現在又過了一個禮拜，還沒收到通知，希望接下來一切順利，趕快把這事了了...續\n\n處理好了，感覺他們有在進步，很樂見Apple在台的第一家直營店愈來愈好，祝他們鴻圖大展！"}]},{"Id":"ChIJzzik3bmrQjQR8TAjh0q04Jk","Name":"OK Mart台北市府門市","Rate":5,"Addr":"110台灣台北市信義區市府路1號B2","Reviews":[{"StoreId":"ChIJzzik3bmrQjQR8TAjh0q04Jk","Text":"注意便利商店在B2，有提供座位。要外帶咖啡需自備杯子。"}]},{"Id":"ChIJJ99mnrCrQjQRBXyAuh93v00","Name":"無印良品 (Muji)","Rate":4.3,"Addr":"110台灣台北市信義區松壽路11號3F","Reviews":[{"StoreId":"ChIJJ99mnrCrQjQRBXyAuh93v00","Text":"咖啡廳的東西比想像中的好吃，價位還可接受，但份量蠻少的較適合稍微填一下肚子。吃完可直接走到旁邊逛書店，雖然藏書不多，但選書感覺得出用心。"},{"StoreId":"ChIJJ99mnrCrQjQRBXyAuh93v00","Text":"Like all Muji stores, this one is bright, tidy, and well-stocked. The zen music playing in the background and the friendly attitude of the service staff make for an overall pleasant shopping experience."},{"StoreId":"ChIJJ99mnrCrQjQRBXyAuh93v00","Text":"Drinks are okay. \nI love the book collections here. \n\n無糖的飲料選擇不多。\n喜歡這裡提供的書籍雜誌閱讀選擇！\n人不多，適合安靜的做些自己想做的事。"},{"StoreId":"ChIJJ99mnrCrQjQRBXyAuh93v00","Text":"空間寬廣，廚具較少。服務親切，附設CAFE，裝潢設計精美，餐食健康養生"},{"StoreId":"ChIJJ99mnrCrQjQRBXyAuh93v00","Text":"由於倉庫遠，店員取貨約要20～30分鐘不等，很悠閒時才能在這裡買東西。"}]},{"Id":"ChIJSTLZ6barQjQRDvsSiktCpo8","Name":"STAY by Yannick Alléno","Rate":4,"Addr":"110台灣台北市信義區市府路45號臺北101購物中心4樓","Reviews":[{"StoreId":"ChIJSTLZ6barQjQRDvsSiktCpo8","Text":"點了NTD1,680套餐：有好吃的餐前麵包，熱前菜推薦香檳慢煮干貝；爐烤羊排及鳀魚脆皮非常好吃，配上菠菜泥(含香菜)與甘藷葉捲非常耳目一新。冷前菜生鮮地中海紅蝦也還不錯。甜點烏龍茶舒芙蕾、烘烤芒果及蛋白霜、爐烤柑橘慕斯及巧克力比起前面的鹹食少了一些層次，下次來還是希望以鹹食為主lol\n\n整個餐廳有好的氛圍，親切的服務還有每一道菜完整的解說，非常好的飲食體驗~"},{"StoreId":"ChIJSTLZ6barQjQRDvsSiktCpo8","Text":"Good food, good service. You can tell that fresh and quality ingredients used  when you eat each and every dish. We had a dinner set. Some of the options in set menu requires additional price. We picked some of them. And it was definetely worth. A small anniversary celebration cake was a nice touch."},{"StoreId":"ChIJSTLZ6barQjQRDvsSiktCpo8","Text":"Great food, as well as good service. Foods are very fresh, delicate and most important of all, delicious! Highly recommend for the set menu which the portion is pretty big. Definitely enough for a big guy, and a little bit too much for a woman."},{"StoreId":"ChIJSTLZ6barQjQRDvsSiktCpo8","Text":"基本上每年都會造訪一次的餐廳，很可惜即將於9月底結束營業。個人很喜歡這邊的氛圍，餐點口味精巧細緻。台北從此少了一個好的法式料理餐廳"},{"StoreId":"ChIJSTLZ6barQjQRDvsSiktCpo8","Text":"餐點份量沒有想像中那麼足而且牛骨髓腥味太重另人難以下嚥，服務人員回達問題有待加強，問是否可加麵包？竟回答：要看廚房還有沒有“剩”，另人傻眼"}]},{"Id":"ChIJmwqDE7qrQjQRI2K5gxuAbck","Name":"春水堂信義店","Rate":3.8,"Addr":"110台灣台北市信義區松壽路9號新光三越信義A9B1","Reviews":[{"StoreId":"ChIJmwqDE7qrQjQRI2K5gxuAbck","Text":"第一次來不知道帶位台在哪，標示不清楚。"},{"StoreId":"ChIJmwqDE7qrQjQRI2K5gxuAbck","Text":"品質還蠻穩定，每間店差異不大。$290栗子燒雞飯，有點像三杯的感覺，不太鹹,配菜有麻婆豆腐(粉辣,有麻的口感)'冷盤花椰菜。"},{"StoreId":"ChIJmwqDE7qrQjQRI2K5gxuAbck","Text":"外帶慢到爆，還一直給別人點餐，是要讓人等多久。外帶點晚餐要去另一個地方取餐也是很不合邏輯"},{"StoreId":"ChIJmwqDE7qrQjQRI2K5gxuAbck","Text":"ร้านชานมไข่มุขชื่อดัง อยู่ชั้น B1 หาไม่ยาก มีโต๊ะ สำหรับนั่งทานที่ร้าน แต่จำนวนน้อยไปหน่อย สามารถสั่งกลับบ้านได้ รสชาติอร่อย ไข่มุขนุ่ม อร่อยคุ้มค่ากับการรอคอย สั่งกลับบ้านอาจจะต้องรอนานหน่อยหากคิวในร้านเยอะ"},{"StoreId":"ChIJmwqDE7qrQjQRI2K5gxuAbck","Text":"珍奶好喝、滷味跟功夫麵都好吃\n份量有點少，價位高\n適合偶爾享受一下~\n若是來台觀光遊客，會推薦來春水堂!"}]},{"Id":"ChIJewvXULqrQjQRZ4KcKM5lPTE","Name":"台北寒舍艾美酒店","Rate":4.3,"Addr":"38 SongRen Road, Xinyi District, Taipei, 台北市台灣 110","Reviews":[{"StoreId":"ChIJewvXULqrQjQRZ4KcKM5lPTE","Text":"CP值不高"},{"StoreId":"ChIJewvXULqrQjQRZ4KcKM5lPTE","Text":":-* "},{"StoreId":"ChIJewvXULqrQjQRZ4KcKM5lPTE","Text":"卓越套房特別舒適，床枕非常舒服，商務設備一應俱全，房間內商務硬體設備非常完善，沐浴品非常舒服，還跟飯店加購，沒有販售部，但可以購買真的很棒。行政貴賓廳倒是非常不完善，電腦得到商務中心用（醬怎麼稱“行政”貴賓廳？）。軟體服務這塊很不一致，早餐人少沒人坐的沙發不給坐，硬是說有訂位？？？第一次聽到自助早餐可以訂位（笑）。但第二天人超多反而可以坐沙發。飯店人員實在也不知道在勢利眼神馬，沒有笑容真的不打緊，夠專業就好，不知道薯條是可以用手拿的嗎？整體服務品質專業度挺落漆的。Room service挺好的，就是…………蠻陽春的餐盤（商務到是ok）跟像印表機印出來的菜單字又小，集點卡式菜單商務飯店沒問題的，但是艾美只是商務的話行政貴賓廳的設備又非常不完整？？？所以定位很怪，所以難怪服務人員態度也混亂吧。但是有幾樣餐點比附近百貨公司有特色多了，所以都在飯店內用餐就好，出門超熱的。"},{"StoreId":"ChIJewvXULqrQjQRZ4KcKM5lPTE","Text":"櫃檯人員跟當初說好的完全不一樣\n我們三人訂兩人房 然後因為床太小後來升級成比較好的房間 那位服務人員說升級有加一張床跟一客早餐 \n但後來我們去吃早餐時只有兩份 去跟櫃檯協調也說沒有這樣說 跟當初商量的不一樣 讓我們多付979\n下次不會再來了"},{"StoreId":"ChIJewvXULqrQjQRZ4KcKM5lPTE","Text":"再過六週 #中秋佳節🎑 即將到來\n在思考挑選什麼好禮與親友分享嗎？\n\n今年 #寒舍集團 推出【#璀艷•#翫月】三款「寒舍中秋月餅禮盒」系列，其中包含「#賞月」廣式月餅禮盒(六入)、「#饌月」典藏月餅禮盒(八入)、「#品月」蛋黃酥禮盒(八入)。\n\n今天來開箱廣式月餅「饌月」禮盒，#抹茶甘栗 、#奶黃咖哩 、#棗泥核桃 、#蓮蓉蛋黃，共四款口味各兩入。\n\nAndy不太愛吃廣式月餅，竟覺得還不錯，每款口味都吃吃看，不過熱量高，只能淺嚐即止！\n\n【這不是葉佩雯 只是開箱文】\n\n#中秋節 #月餅 #月 #十五夜 #寒舍 #艾美 #寒舍艾美 #喜來登 #台北 #taipei\n\n寒舍艾美酒店 台北寒舍艾美酒店 Le Méridien Taipei 台北寒舍艾美酒店"}]},{"Id":"ChIJK8sU3LarQjQR0WE6XpCrPKI","Name":"PANDORA 台北101概念店","Rate":3,"Addr":"110台灣台北市信義區市府路45號","Reviews":[{"StoreId":"ChIJK8sU3LarQjQR0WE6XpCrPKI","Text":"如果你也遇到店員懶散一付不做台灣年輕人的態度，那走去新光三越A11會是個好選擇！"},{"StoreId":"ChIJK8sU3LarQjQR0WE6XpCrPKI","Text":"GOOD"}]},{"Id":"ChIJBeS-dLqrQjQRYTLib85vH4I","Name":"PANDORA 新光三越 A11","Rate":0,"Addr":"110台灣台北市信義區松壽路11號1 樓","Reviews":[]},{"Id":"ChIJkXelC7qrQjQRHIwaUKR9qUY","Name":"Vivienne Westwood (新光三越A9)","Rate":0,"Addr":"11051台灣台北市信义区松壽路9號3樓","Reviews":[]},{"Id":"ChIJ0_P1N7qrQjQRIdfS6R5-ht8","Name":"HAAGEN DAZS 哈根達斯","Rate":4.5,"Addr":"110台灣台北市信義區松高路12號B2","Reviews":[{"StoreId":"ChIJ0_P1N7qrQjQRIdfS6R5-ht8","Text":"座位雖少但值得等待，可悠閒的享用美味冰淇淋下午茶與咖啡，坐看過往行人百態。"},{"StoreId":"ChIJ0_P1N7qrQjQRIdfS6R5-ht8","Text":"知名冰淇淋的門市，各式不同種類的冰淇淋，也與一般只銷售的門市不同，有提供休息的沙發位位區，菜單除單品外也有甜點組合，更有少見的咖啡，內用外帶皆可，不錯吃。"},{"StoreId":"ChIJ0_P1N7qrQjQRIdfS6R5-ht8","Text":"冰淇淋很好吃"},{"StoreId":"ChIJ0_P1N7qrQjQRIdfS6R5-ht8","Text":"抹茶"},{"StoreId":"ChIJ0_P1N7qrQjQRIdfS6R5-ht8","Text":""}]},{"Id":"ChIJSTLZ6barQjQRe92tH5f2atI","Name":"H:CONNECT 台北101旗艦店","Rate":4,"Addr":"110台灣台北市信義區市府路45號B1","Reviews":[{"StoreId":"ChIJSTLZ6barQjQRe92tH5f2atI","Text":""}]},{"Id":"ChIJcxkKdLqrQjQRXICV9IUkFm0","Name":"H:CONNECT 新光三越A11專櫃","Rate":5,"Addr":"110台灣台北市信义区松壽路11號號 B1","Reviews":[{"StoreId":"ChIJcxkKdLqrQjQRXICV9IUkFm0","Text":"服務生服務很好 很會介紹！ 這家很好買！"}]},{"Id":"ChIJG2UmvbCrQjQRSVv9nvlDR0Q","Name":"H:CONNECT 信義威秀旗艦店","Rate":4.7,"Addr":"110台灣台北市信義區松壽路18號","Reviews":[{"StoreId":"ChIJG2UmvbCrQjQRSVv9nvlDR0Q","Text":"Good quality clothing with nice designs at a reasonable price. "},{"StoreId":"ChIJG2UmvbCrQjQRSVv9nvlDR0Q","Text":"有少女時代 Yoona 廣告的女裝店"},{"StoreId":"ChIJG2UmvbCrQjQRSVv9nvlDR0Q","Text":"GOOD"}]},{"Id":"ChIJy1N6J7erQjQRaOn7RHYMfJY","Name":"Scotch \u0026 Soda","Rate":3,"Addr":"No. 7, Section 5, Xinyi Rd, 信義區台北市台灣 110","Reviews":[{"StoreId":"ChIJy1N6J7erQjQRaOn7RHYMfJY","Text":"Young fashion that appeals to an older age group"}]},{"Id":"ChIJXVRz3LmrQjQRZx-tF_v1E2U","Name":"臺北市政府","Rate":3.8,"Addr":"110台灣台北市信義區市府路1號","Reviews":[{"StoreId":"ChIJXVRz3LmrQjQRZx-tF_v1E2U","Text":"建築物量體很大,年代有點久,當時是因為雙十平面設計與國慶日有關競圖時被選上,有許多台北市的公務機關合署辦公,如果不知道要找甚麼單位處理問題,可以先到一樓後側有一個聯合服務中心洽詢,地下室還有便宜的自助餐可以吃午飯"},{"StoreId":"ChIJXVRz3LmrQjQRZx-tF_v1E2U","Text":"台北市政府一直是台北人心中重要地標之一！也是台北市裡辦公的好地方！自從台北市政府換了柯市長以後！是正更是進步！內部人文薈萃！還有用餐的餐廳！還有咖啡廳！更常常不時地在中庭舉辦文藝活動！實在是個辦公的好地方！也是我們台北市民的驕傲！！！"},{"StoreId":"ChIJXVRz3LmrQjQRZx-tF_v1E2U","Text":"地下b2常有不同的廠商進駐，有時候可以挖到寶喔"},{"StoreId":"ChIJXVRz3LmrQjQRZx-tF_v1E2U","Text":"樓下B2 有攤位 和便利商店，美容院，送洗，修鞋，刻印章，很便利"},{"StoreId":"ChIJXVRz3LmrQjQRZx-tF_v1E2U","Text":"各區公所重要市政業務都要來此辦理……"}]},{"Id":"ChIJhRzNZrarQjQRi4np_EimKkQ","Name":"中美洲經貿辦事處","Rate":0,"Addr":"110台灣信義區信義路五段5號","Reviews":[]},{"Id":"ChIJA-L9zLGrQjQRBiyvB1Zhu8o","Name":"捷運象山站","Rate":0,"Addr":"110台灣台北市信義區","Reviews":[]},{"Id":"ChIJDYIXhLCrQjQREd9BK4ylpXM","Name":"LAVA","Rate":3.4,"Addr":"110台灣台北市信義區松壽路22號B1","Reviews":[{"StoreId":"ChIJDYIXhLCrQjQREd9BK4ylpXM","Text":"\n\n"},{"StoreId":"ChIJDYIXhLCrQjQREd9BK4ylpXM","Text":"很爛耶 要不是朋友找實在不想來這\n看個證件久就算了 音樂也跟不上其他間酒也是\n再來 我第一次聽到台北夜店請一個聽都沒聽過的團體來唱歌 還唱中文？？？\n門票還敢收那麼貴😆近來聽中文歌的耶"},{"StoreId":"ChIJDYIXhLCrQjQREd9BK4ylpXM","Text":"音樂跟不上主流，DJ話很多\n酒有點淡，但仍然可以盡興"},{"StoreId":"ChIJDYIXhLCrQjQREd9BK4ylpXM","Text":"酒出的快，公關，安管，外場，吧檯人都很和善唷！！！"},{"StoreId":"ChIJDYIXhLCrQjQREd9BK4ylpXM","Text":"空調太差讓人快呼吸不到空氣應該立即改善。"}]},{"Id":"ChIJfdvfrbCrQjQRzw0XcXO0gJ4","Name":"COMMUNE A7","Rate":4.2,"Addr":"110台灣台北市信義區松壽路3號","Reviews":[{"StoreId":"ChIJfdvfrbCrQjQRzw0XcXO0gJ4","Text":"很有特色的貨櫃，常常會辦一些活動很有趣"},{"StoreId":"ChIJfdvfrbCrQjQRzw0XcXO0gJ4","Text":"裡面很多好吃的店🤗🤗"},{"StoreId":"ChIJfdvfrbCrQjQRzw0XcXO0gJ4","Text":"很好拍照😍"},{"StoreId":"ChIJfdvfrbCrQjQRzw0XcXO0gJ4","Text":"💝💝💝"},{"StoreId":"ChIJfdvfrbCrQjQRzw0XcXO0gJ4","Text":"2017/08/19 世大運開幕第一天，非常棒的露天環境，恣意享受輕鬆愉快的氛圍，涵蓋吃喝玩樂、欣賞表演，打造出都市悠閒風，如同置身於國外假日市集，走走逛逛到處充滿驚喜，為於信義繁華商圈，熱鬧至始至終感染人心。"}]}]`
	fmt.Fprint(w, mockData)
}

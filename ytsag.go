package main

import (
	 "encoding/json"
	 "io/ioutil"
	 "net/http"
	 "net/url"
	 "regexp"
	 "errors"
	 "strconv"
	 "strings"
)

type CmdConfigYts struct {
	 Reg_rand        *regexp.Regexp
	 Reg_detail      *regexp.Regexp
	 Reg_search        *regexp.Regexp
}

var Yts_randof int

func NewCmdYts(config *CmdConfigYts) {
	 config.Reg_rand = regexp.MustCompile(`^/random_yts(?:(@[a-zA-Z0-9_]+bot)?( [0-9]+)?$)`)
	 config.Reg_detail = regexp.MustCompile(`^/details_yts_[0-9]+$`)
	 config.Reg_search = regexp.MustCompile(`^/search_yts(?:(@[a-zA-Z0-9_]+bot)? .+$)`)
	 Yts_randof = 1
}

func YtsSetPage(a string){
	 b := strings.SplitN(a, " ", 2)
	 if len(b) == 2{
			Yts_randof, _ = strconv.Atoi(b[1])
	 }else{
			Yts_randof = 1
	 }
}

func YtsRandom(gen string) (string, error){
	 type movie struct{
			ID int   `json:"id"`
			Title string   `json:"title_long"`
			Rating   float32  `json:"rating"`
			Genres   []string `json:"genres"`
	 }
	 type t_data struct{
			Movie_count int   `json:"movie_count"`
			Movies   []movie  `json:"movies"`
	 }
	 var respData struct {
			Data t_data `json:"data"`
	 }

	 limit := 5
	 qResp, err := YtsQuery(limit, "", Yts_randof, "1080p", "4", gen)
	 if err != nil{
			return "", err
	 }
	 err = json.Unmarshal(qResp, &respData)
	 if err != nil {
			return "", err
	 }
	 if respData.Data.Movie_count == 0{
			return "No movies found!", nil
	 }
	 if gen == ""{gen = "All"}
	 m := "Page: " + strconv.Itoa(Yts_randof) + " - Category: " + gen + "\n"
	 if Yts_randof*limit >= respData.Data.Movie_count{
			Yts_randof = 1
	 }else{
			Yts_randof++
	 }

	 for _, i := range respData.Data.Movies{
			s_genres := ": "
			for _, j := range i.Genres{
				 s_genres += j + ", "
			}
			m += "\nTítulo: " + i.Title + "\nRating: " +
			strconv.FormatFloat(float64(i.Rating), 'f', -1, 32) + "\nGenres" +
			s_genres[:len(s_genres)-2] + "\nDetails and Download: /details_yts_" + strconv.Itoa(i.ID) + "\n"
	 }
	 return m, nil
}

func YtsSearch(a string) (string, error) {
	 type movie struct{
			ID int   `json:"id"`
			Title string   `json:"title_long"`
			Rating   float32  `json:"rating"`
			Genres   []string `json:"genres"`
	 }
	 type t_data struct{
			Movie_count int   `json:"movie_count"`
			Movies   []movie  `json:"movies"`
	 }
	 var respData struct {
			Data t_data `json:"data"`
	 }

	 b := strings.SplitN(a, " ", 2)
	 if len(b) != 2{
			return "", errors.New("Error in search function")
	 }
	 qResp, err := YtsQuery(50, b[1], 1, "", "", "")
	 if err != nil{
			return "", err
	 }
	 err = json.Unmarshal(qResp, &respData)
	 if err != nil {
			return "", err
	 }
	 if respData.Data.Movie_count == 0{
			return "No movies found!", nil
	 }
	 m := ""
	 for _, i := range respData.Data.Movies{
			s_genres := ": "
			for _, j := range i.Genres{
				 s_genres += j + ", "
			}
			m += "\nTítulo: " + i.Title + "\nRating: " +
			strconv.FormatFloat(float64(i.Rating), 'f', -1, 32) + "\nGenres" +
			s_genres[:len(s_genres)-2] + "\nDetails and Download: /details_yts_" + strconv.Itoa(i.ID) + "\n"
	 }
	 return m, nil
}

func YtsDetail(a string) (string,string,string,error){
	 type torrent struct{
			Hash   string  `json:"hash"`
	 }
	 type movie struct{
			ID    int      `json:"id"`
			Description   string `json:"description_full"`
			Rating float32  `json:"rating"`
			Title string   `json:"title_long"`
			Img   string   `json:"large_cover_image"`
			Trailer  string   `json:"yt_trailer_code"`
			Torrents  []torrent   `json:"torrents"`
			Genres    []string `json:"genres"`
	 }
	 type t_data struct{
			Movie   movie  `json:"movie"`
	 }
	 var respData struct {
			Data t_data `json:"data"`
	 }
	 b := strings.SplitN(a, "_", 3)
	 var id int
	 var err error
	 if len(b) == 3{
			id, err = strconv.Atoi(b[2])
			if err != nil{
				 return "","", "", err
			}
	 }else{
			return "", "", "", errors.New("Error parsing the command")
	 }
	 resp, err := http.Get("https://yts.ag/api/v2/movie_details.json?movie_id=" + strconv.Itoa(id))
	 if err != nil {
			return "","", "",err
	 }
	 defer resp.Body.Close()
	 if resp.StatusCode != http.StatusOK {
			return "","", "", errors.New("HTTP Status: " + strconv.Itoa(resp.StatusCode))
	 }
	 repBody, err := ioutil.ReadAll(resp.Body)
	 if err != nil {
			return "","", "", err
	 }
	 err = json.Unmarshal(repBody, &respData)
	 if err != nil {
			return "","", "", err
	 }
	 if respData.Data.Movie.ID == 0{
			return "", "No movie found with this ID", "", nil
	 }
	 var s_torrent string
	 if len(respData.Data.Movie.Torrents) > 1{
			s_torrent = respData.Data.Movie.Torrents[1].Hash
	 }else{
			s_torrent = respData.Data.Movie.Torrents[0].Hash
	 }
	 s_genres := ": "
	 for _, i := range respData.Data.Movie.Genres{
			s_genres += i + ", "
	 }
	 m := "+ Título: " + respData.Data.Movie.Title + "\n+ Descripción: " +
	 respData.Data.Movie.Description + "\n+ Géneros" + s_genres[:len(s_genres)-2] + "\n+ Rating: " +
	 strconv.FormatFloat(float64(respData.Data.Movie.Rating), 'f', -1, 32) +
	 "\n+ Trailer: " + "https://www.youtube.com/watch?v=" + respData.Data.Movie.Trailer
	 return respData.Data.Movie.Img, m, s_torrent, nil
}


func YtsQuery (limit int, term string, page int, quality string, rating string, genre string) ([]byte, error){
	 resp, err := http.Get("https://yts.ag/api/v2/list_movies.json?limit=" + strconv.Itoa(limit) +
	 "&page=" + strconv.Itoa(page) + "&query_term=" + url.QueryEscape(term) + "&quality=" + quality +
	 "&minimum_rating=" + rating + "&genre=" + genre + "&sort_by=year&order_by=desc")

	 if err != nil {
			return []byte{}, err
	 }
	 defer resp.Body.Close()
	 if resp.StatusCode != http.StatusOK {
			return []byte{}, errors.New("HTTP Status: " + strconv.Itoa(resp.StatusCode))
	 }
	 repBody, err := ioutil.ReadAll(resp.Body)
	 if err != nil {
			return []byte{}, err
	 }
	 return repBody, nil
}

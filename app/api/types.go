package api

import "encoding/xml"

type memberValue struct {
	Text    string `xml:",chardata"`
	Boolean string `xml:"boolean"`
	String  string `xml:"string"`
	//DateTimeIso8601 string `xml:"dateTime.iso8601"`
	//Array           array  `xml:"array"`
	//Int             string `xml:"int"`
}

type member struct {
	Text  string      `xml:",chardata"`
	Name  string      `xml:"name"`
	Value memberValue `xml:"value"`
}

type valStruct struct {
	Text   string   `xml:",chardata"`
	Member []member `xml:"member"`
}

type dataValue struct {
	Text   string    `xml:",chardata"`
	String string    `xml:"string"`
	Struct valStruct `xml:"struct"`
}

type data struct {
	Text  string    `xml:",chardata"`
	Value dataValue `xml:"value"`
}

type array struct {
	Text string `xml:",chardata"`
	Data data   `xml:"data"`
}
type value struct {
	Text  string `xml:",chardata"`
	Array array  `xml:"array"`
}

type param struct {
	Text  string `xml:",chardata"`
	Value value  `xml:"value"`

	//DateTimeIso8601 string `xml:"dateTime.iso8601"`
	//String          string `xml:"string"`
}

type params struct {
	Text  string `xml:",chardata"`
	Param param  `xml:"param"`
}

type getUsersBlogsResponse struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  params   `xml:"params"`
}

type getCategoriesResponse struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  struct {
		Text  string `xml:",chardata"`
		Param struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text  string `xml:",chardata"`
				Array struct {
					Text string `xml:",chardata"`
					Data struct {
						Text  string      `xml:",chardata"`
						Value []dataValue `xml:"value"`
					} `xml:"data"`
				} `xml:"array"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

type getPostResponse struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  struct {
		Text  string `xml:",chardata"`
		Param struct {
			Text  string    `xml:",chardata"`
			Value dataValue `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

type postRequest struct {
	XMLName    xml.Name `xml:"methodCall"`
	Text       string   `xml:",chardata"`
	MethodName string   `xml:"methodName"`
	Params     struct {
		Text  string `xml:",chardata"`
		Param []struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text   string `xml:",chardata"`
				String string `xml:"string"`
				Struct struct {
					Text   string `xml:",chardata"`
					Member []struct {
						Text  string `xml:",chardata"`
						Name  string `xml:"name"`
						Value struct {
							Text    string `xml:",chardata"`
							String  string `xml:"string"`
							Boolean string `xml:"boolean"`
							Base64  string `xml:"base64"`
							Array   struct {
								Text string `xml:",chardata"`
								Data []struct {
									Text  string `xml:",chardata"`
									Value struct {
										Text   string `xml:",chardata"`
										String string `xml:"string"`
									} `xml:"value"`
								} `xml:"data"`
							} `xml:"array"`
							DateTimeIso8601 string `xml:"dateTime.iso8601"`
						} `xml:"value"`
					} `xml:"member"`
				} `xml:"struct"`
				Boolean string `xml:"boolean"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

type mediaRequest struct {
	XMLName    xml.Name `xml:"methodCall"`
	Text       string   `xml:",chardata"`
	MethodName string   `xml:"methodName"`
	Params     struct {
		Text  string `xml:",chardata"`
		Param []struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text   string `xml:",chardata"`
				String string `xml:"string"`
				Struct struct {
					Text   string `xml:",chardata"`
					Member []struct {
						Text  string `xml:",chardata"`
						Name  string `xml:"name"`
						Value struct {
							Text    string `xml:",chardata"`
							Boolean string `xml:"boolean"`
							Base64  string `xml:"base64"`
							String  string `xml:"string"`
						} `xml:"value"`
					} `xml:"member"`
				} `xml:"struct"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

type newPostResponse struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  struct {
		Text  string `xml:",chardata"`
		Param struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text   string `xml:",chardata"`
				String string `xml:"string"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

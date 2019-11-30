package controller

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/websocket"
	"github.com/yanzhen74/goview/src/model"
)

// backend frame model
var Dicts *[]model.FrameDict

// map from para index of frame to para index of view

func Init_0c_Processer(conf string) {
	z, _ := model.Read_para_dict(conf)
	Dicts = model.Get_frame_dict_list(z)

	for _, d := range *Dicts {
		Bind_network(d.Frame_type)
		go Process0cPkg(d)
	}
}

func Process0cPkg(frame model.FrameDict) {
	para_view_map := make(map[*websocket.NSConn]map[int]string)

	// pkg should send only required parameters to view's chan
	var pkg map[chan string]interface{} = make(map[chan string]interface{})

	ticker := time.NewTicker(time.Millisecond * time.Duration(100))
	cases := init_cases(frame.Frame_type.NetChanFrame, ticker, frame.Frame_type.UserChanReg)

	e := 0
	for i := 0; ; {
		chose, value, _ := reflect.Select(cases)

		// log.Printf("cases len %d, channel no %s %d %d\n", len(cases), frame.Frame_type.MissionID, len(frame.Frame_type.UserChanMap), i)
		switch chose {
		case 0: // regist/unregist chan_view
			info := (value.Interface().(*model.View_page_regist_info))
			if -1 == regist_view_chan(&frame, info, para_view_map) {
				delete(pkg, info.View_chan)
				// todo should remove this view_chan from cases
				cases = cases[:3]
			}
		case 1: // time
			// to be deleted, simulate net receiver
			//frame.Frame_type.NetChanFrame <- "hello world"
		case 2: // net frame
			// update when receive net data
			msg := (value.Interface()).(string)

			// todo zcy do this function
			v, err := get_param_array_from_frame(i, &frame, msg)
			if err != nil {
				e++
				continue
			}

			var buffer bytes.Buffer
			for conn, view_chan := range frame.Frame_type.UserChanMap {
				buffer.Reset()
				for id_in_frame, id_in_view := range para_view_map[conn] {
					// if param not in this frame
					if value, ok := v[id_in_frame]; ok {
						buffer.WriteString(id_in_view)
						buffer.WriteString(value)
					}
				}
				pkg[view_chan] = buffer.String()
			}
			// ? when one view is blocked, should not sent to it
			// todo : should delete the block one, not tails
			if len(cases) > 3+len(frame.Frame_type.UserChanMap) {
				cases = cases[:len(cases)-len(frame.Frame_type.UserChanMap)]
			}
			send_to_view(&cases, frame.Frame_type.UserChanMap, pkg)
			i++
		default:
			cases = append(cases[:chose], cases[chose+1:]...)
		}
	}
}

// zcy do this
// input hello:
// #<DATA_TYPE>_<MISSIONID>_<SubAddress>_<DataFormat>\t<station>\t<FrameNo>\t<GroundTime>
// RTM_TGTH_PK-CEH2_Result  00  0   00:00:03.3345
// #<DataItemID> <DataItemCode> <DataItemResult>;...
// 1 00000000 0.000;2 00000000 0.000;33 ee 238;
// output v:
// v[j]:",<DataItemCode>,<DataItemResult>,<Description>,<-1:out of limit;0:normal>;";
// ,00000000,0.000,正常,0;,00000000,0.000,异常,-1;,ee,238,正确,0;
func get_param_array_from_frame(i int, frame *model.FrameDict, msg string) (v map[int]string, err error) {
	err = nil
	v = make(map[int]string)

	lines := strings.Split(msg, "\n")
	if len(lines) <= 1 {
		err = errors.New("msg has no param line")
		return nil, err
	}
	params := strings.Split(lines[1], ";")
	if len(params) <= 0 {
		err = errors.New("msg has no params")
		return nil, err
	}

	j := 0
	for _, x := range params {
		p := strings.Split(x, " ")
		if len(p) < 3 {
			continue
		}

		// p[0]:id, p[1]: code, p[2]: result
		id, _ := strconv.Atoi(p[0])
		for j < len(frame.ParaList) && frame.ParaList[j].ID < id {
			j++
		}
		for j < len(frame.ParaList) && frame.ParaList[j].ID == id {
			// begin --------------------
			// zcy do this section
			// from p to v[j]
			// this is just for test; for example: 0 should be -1 if out of limit
			strCode := &p[1]
			Temp := p[2]
			strResult := &Temp
			strResultValue := &p[2]
			Normal, Error := Param_Transfer(frame.ParaList[j], strCode, strResult, strResultValue)
			v[j] = fmt.Sprintf(",0x%s,%s,%s,%d;", p[1], p[2], *strResult, Normal)
			if Error != nil {
				fmt.Printf("\n%s\n", Error)
			}
			// end ------------------------
			j++
		}
	}

	return v, err
}

func Param_Transfer(para model.Para, strCode *string, strResult *string, strResultValue *string) (Normal int, Error error) {
	Error = nil
	process_type := para.Process_type
	if strings.ToLower(process_type) == "" {
		*strCode = *strCode
		*strResultValue = *strCode
		Normal = 0
		return Normal, Error
	} else if strings.ToLower(process_type) == "raw" {
		process_unit := para.Process_unit
		process_start := para.Process_start
		process_end := para.Process_end

		if strings.ToLower(process_unit) == "byte" {
			*strResult = *strCode
			posStart, errstart := strconv.ParseInt(process_start, 10, 32)
			endStart, errend := strconv.ParseInt(process_end, 10, 32)
			length := endStart - posStart + 1
			proCode := *strCode

			if errstart != nil || errend != nil || int(length) > len(proCode) {
				if int(length) > len(proCode) {
					Normal = -1
					Error = errors.New("Fatal:Limitation Break in Byte\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
				} else {
					Normal = -1
					Error = errors.New("Error:PosStart or EndStart Error in Byte\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
				}
				return Normal, Error
			}

			proCode = proCode[posStart : posStart+length] // posStart从零还是从1开始？

			for k := 0; k < len(para.ParaRangeList); k++ {
				para.ParaRangeList[k].Alarm_max_equal = strings.ToLower(para.ParaRangeList[k].Alarm_max_equal)
				para.ParaRangeList[k].Alarm_min_equal = strings.ToLower(para.ParaRangeList[k].Alarm_min_equal)
				if strings.ToLower(para.ParaRangeList[k].Alarm_max) == strings.ToLower(para.ParaRangeList[k].Alarm_min) {
					if strings.ToLower(para.ParaRangeList[k].Alarm_min) == strings.ToLower(proCode) {
						if strings.ToLower(para.ParaRangeList[k].ParaRangeSpecification) != "" {
							*strResult = para.ParaRangeList[k].ParaRangeSpecification
							*strCode = *strCode
							*strResultValue = proCode
							Normal = 0
							return Normal, Error
						} else {
							*strResult = proCode
							*strCode = *strCode
							*strResultValue = proCode
							Normal = 0
							return Normal, Error
						}
					}
				} else {
					sup, errsup := strconv.ParseInt(para.ParaRangeList[k].Alarm_max, 16, 64)
					inf, errinf := strconv.ParseInt(para.ParaRangeList[k].Alarm_min, 16, 64)
					value, errval := strconv.ParseInt(proCode, 16, 64)

					if errsup != nil || errinf != nil || errval != nil {
						Normal = -1
						Error = errors.New("Error:Sup, Inf or Val Error in Bit\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
						return Normal, Error
					}

					if para.ParaRangeList[k].Alarm_max_equal == "true" && para.ParaRangeList[k].Alarm_min_equal == "true" {
						if value >= inf && value <= sup {
							if para.ParaRangeList[k].ParaRangeSpecification != "" {
								*strResult = para.ParaRangeList[k].ParaRangeSpecification
								*strCode = *strCode
								*strResultValue = proCode
								Normal = 0
								return Normal, Error
							} else {
								if para.Type == "int" {
									*strResult = strconv.FormatInt(value, 10)
								} else {
									*strResult = proCode
								}
								*strCode = *strCode
								*strResultValue = proCode
								Normal = 0
								return Normal, Error
							}
						}
					} else if para.ParaRangeList[k].Alarm_max_equal == "false" && para.ParaRangeList[k].Alarm_min_equal == "true" {
						if value >= inf && value < sup {
							if para.ParaRangeList[k].ParaRangeSpecification != "" {
								*strResult = para.ParaRangeList[k].ParaRangeSpecification
								*strCode = *strCode
								*strResultValue = proCode
								Normal = 0
								return Normal, Error
							} else {
								if para.Type == "int" {
									*strResult = strconv.FormatInt(value, 10)
								} else {
									*strResult = proCode
								}
								*strCode = *strCode
								*strResultValue = proCode
								Normal = 0
								return Normal, Error
							}
						}
					} else if para.ParaRangeList[k].Alarm_max_equal == "true" && para.ParaRangeList[k].Alarm_min_equal == "false" {
						if value > inf && value <= sup {
							if para.ParaRangeList[k].ParaRangeSpecification != "" {
								*strResult = para.ParaRangeList[k].ParaRangeSpecification
								*strCode = *strCode
								*strResultValue = proCode
								Normal = 0
								return Normal, Error
							} else {
								if para.Type == "int" {
									*strResult = strconv.FormatInt(value, 10)
								} else {
									*strResult = proCode
								}
								*strCode = *strCode
								*strResultValue = proCode
								Normal = 0
								return Normal, Error
							}
						}
					} else if para.ParaRangeList[k].Alarm_max_equal == "false" && para.ParaRangeList[k].Alarm_min_equal == "false" {
						if value > inf && value < sup {
							if para.ParaRangeList[k].ParaRangeSpecification != "" {
								*strResult = para.ParaRangeList[k].ParaRangeSpecification
								*strCode = *strCode
								*strResultValue = proCode
								Normal = 0
								return Normal, Error
							} else {
								if para.Type == "int" {
									*strResult = strconv.FormatInt(value, 10)
								} else {
									*strResult = proCode
								}
								*strCode = *strCode
								*strResultValue = proCode
								Normal = 0
								return Normal, Error
							}
						}
					}
				}
			}
			*strCode = *strCode
			*strResultValue = proCode
		} else if strings.ToLower(process_unit) == "bit" {
			posStart, errstart := strconv.ParseInt(process_start, 10, 32)
			endStart, errend := strconv.ParseInt(process_end, 10, 32)
			codeLen := len(*strCode)
			length := endStart - posStart + 1
			totalbinary := ""
			if errstart != nil || errend != nil {
				Normal = -1
				Error = errors.New("Error:PosStart or EndStart Error in Bit\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
				return Normal, Error
			}

			for k := 0; k < codeLen; k++ {
				segcode := *strCode
				segcode = segcode[k : k+1]
				intbinacode, errbina := strconv.ParseInt(segcode, 16, 32)
				binacode := strconv.FormatInt(intbinacode, 2)
				_0bina := "0000"
				fmt.Print("")
				if errbina != nil {
					Normal = -1
					Error = errors.New("Error:Intbinacode Error in Bit\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
					return Normal, Error
				}
				if len(binacode) < 4 {
					_0bina := _0bina[0:(4 - len(binacode))]
					binacode = _0bina + binacode
				}
				totalbinary = totalbinary + binacode
			}
			bincode := ""
			for k := 0; k < len(totalbinary); k++ {
				bincode = bincode + "."
			}
			proCode := ""
			proCode = totalbinary[posStart : posStart+length]
			// replaceBinacode := bincode[posStart : posStart+length]
			// sb := bincode
			// sb = strings.Replace(sb, replaceBinacode, proCode, 1)
			sb := bincode[:posStart] + proCode + bincode[posStart+length:]
			// sb = sb[0:posStart] + proCode + sb[posStart+length:len(bincode)]
			*strResult = sb

			v, errv := strconv.ParseInt(proCode, 2, 32)
			if errv != nil {
				Normal = -1
				Error = errors.New("Error:V Error in Bit\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
				return Normal, Error
			}

			for k := 0; k < len(para.ParaRangeList); k++ {
				para.ParaRangeList[k].Alarm_max_equal = strings.ToLower(para.ParaRangeList[k].Alarm_max_equal)
				para.ParaRangeList[k].Alarm_min_equal = strings.ToLower(para.ParaRangeList[k].Alarm_min_equal)
				if strings.ToLower(para.ParaRangeList[k].Alarm_max) == strings.ToLower(para.ParaRangeList[k].Alarm_min) {
					if strings.ToLower(para.ParaRangeList[k].Alarm_min) == strings.ToLower(proCode) {
						if strings.ToLower(para.ParaRangeList[k].ParaRangeSpecification) != "" {
							*strResult = para.ParaRangeList[k].ParaRangeSpecification
							*strCode = *strCode
							if strings.ToLower(para.Type) == "int" {
								*strResultValue = strconv.FormatInt(v, 10)
							} else {
								*strResultValue = sb
							}
						} else {
							if strings.ToLower(para.Type) == "int" {
								*strResult = strconv.FormatInt(v, 10)
							} else {
								*strResult = sb
							}
							*strCode = *strCode
							if strings.ToLower(para.Type) == "int" {
								*strResultValue = strconv.FormatInt(v, 10)
							} else {
								*strResultValue = sb
							}
						}
						Normal = 0
						return Normal, Error
					}
				} else {
					sup, errsup := strconv.ParseInt(para.ParaRangeList[k].Alarm_max, 10, 64)
					inf, errinf := strconv.ParseInt(para.ParaRangeList[k].Alarm_min, 10, 64)
					value, errval := strconv.ParseInt(proCode, 2, 32)

					if errsup != nil || errinf != nil || errval != nil {
						Normal = -1
						Error = errors.New("Error:Sup, Inf or Value Error in Bit\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
						return Normal, Error
					}

					if para.ParaRangeList[k].Alarm_max_equal == "true" && para.ParaRangeList[k].Alarm_min_equal == "true" {
						if value >= inf && value <= sup {
							if para.ParaRangeList[k].ParaRangeSpecification != "" {
								*strResult = para.ParaRangeList[k].ParaRangeSpecification
								*strCode = *strCode
								if strings.ToLower(para.Type) == "int" {
									*strResultValue = strconv.FormatInt(v, 10)
								} else {
									*strResultValue = sb
								}
							} else {
								if strings.ToLower(para.Type) == "int" {
									*strResult = strconv.FormatInt(v, 10)
								} else {
									*strResult = sb
								}
								*strCode = *strCode
								if strings.ToLower(para.Type) == "int" {
									*strResultValue = strconv.FormatInt(v, 10)
								} else {
									*strResultValue = sb
								}
							}
							Normal = 0
							return Normal, Error
						}
					} else if para.ParaRangeList[k].Alarm_max_equal == "false" && para.ParaRangeList[k].Alarm_min_equal == "true" {
						if value >= inf && value < sup {
							if para.ParaRangeList[k].ParaRangeSpecification != "" {
								*strResult = para.ParaRangeList[k].ParaRangeSpecification
								*strCode = *strCode
								if strings.ToLower(para.Type) == "int" {
									*strResultValue = strconv.FormatInt(v, 10)
								} else {
									*strResultValue = sb
								}
							} else {
								if strings.ToLower(para.Type) == "int" {
									*strResult = strconv.FormatInt(v, 10)
								} else {
									*strResult = sb
								}
								*strCode = *strCode
								if strings.ToLower(para.Type) == "int" {
									*strResultValue = strconv.FormatInt(v, 10)
								} else {
									*strResultValue = sb
								}
							}
							Normal = 0
							return Normal, Error
						}
					} else if para.ParaRangeList[k].Alarm_max_equal == "true" && para.ParaRangeList[k].Alarm_min_equal == "false" {
						if value > inf && value <= sup {
							if para.ParaRangeList[k].ParaRangeSpecification != "" {
								*strResult = para.ParaRangeList[k].ParaRangeSpecification
								*strCode = *strCode
								if strings.ToLower(para.Type) == "int" {
									*strResultValue = strconv.FormatInt(v, 10)
								} else {
									*strResultValue = sb
								}
							} else {
								if strings.ToLower(para.Type) == "int" {
									*strResult = strconv.FormatInt(v, 10)
								} else {
									*strResult = sb
								}
								*strCode = *strCode
								if strings.ToLower(para.Type) == "int" {
									*strResultValue = strconv.FormatInt(v, 10)
								} else {
									*strResultValue = sb
								}
							}
							Normal = 0
							return Normal, Error
						}
					} else if para.ParaRangeList[k].Alarm_max_equal == "false" && para.ParaRangeList[k].Alarm_min_equal == "false" {
						if value > inf && value < sup {
							if para.ParaRangeList[k].ParaRangeSpecification != "" {
								*strResult = para.ParaRangeList[k].ParaRangeSpecification
								*strCode = *strCode
								if strings.ToLower(para.Type) == "int" {
									*strResultValue = strconv.FormatInt(v, 10)
								} else {
									*strResultValue = sb
								}
							} else {
								if strings.ToLower(para.Type) == "int" {
									*strResult = strconv.FormatInt(v, 10)
								} else {
									*strResult = sb
								}
								*strCode = *strCode
								if strings.ToLower(para.Type) == "int" {
									*strResultValue = strconv.FormatInt(v, 10)
								} else {
									*strResultValue = sb
								}
							}
							Normal = 0
							return Normal, Error
						}
					}
				}
			}
			*strCode = strings.Replace(*strCode, "0x", "", -1)
			*strCode = *strCode
			*strResultValue = sb
		}
		if strings.ToLower(process_unit) == "code" {
			*strResult = *strCode
			for k := 0; k < len(para.ParaRangeList); k++ {
				posStart, errstrat := strconv.ParseInt(process_start, 10, 32)
				endStart, errend := strconv.ParseInt(process_end, 10, 32)
				length := endStart - posStart + 1
				proCode := *strCode
				proCode = proCode[posStart : posStart+length] // posStart从零还是从1开始？

				if errstrat != nil || errend != nil {
					Normal = -1
					Error = errors.New("Error:PosStrat or EndStart Error in Code\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
					return Normal, Error
				}
				if strings.ToLower(proCode) == strings.ToLower(para.ParaRangeList[k].Alarm_max) && strings.ToLower(proCode) == strings.ToLower(para.ParaRangeList[k].Alarm_min) {
					if strings.ToLower(para.ParaRangeList[k].ParaRangeSpecification) != "" {
						*strResult = para.ParaRangeList[k].ParaRangeSpecification
						*strCode = *strCode
						*strResultValue = proCode
						Normal = 0
						return Normal, Error
					} else {
						*strResult = proCode
						*strCode = *strCode
						*strResultValue = proCode
						Normal = 0
						return Normal, Error
					}
				}
				*strResultValue = proCode
			}
		}
		if strings.ToLower(process_unit) == "longcode" {
			*strResult = *strCode
			for k := 0; k < len(para.ParaRangeList); k++ {
				posStart, errstrat := strconv.ParseInt(process_start, 10, 32)
				endStart, errend := strconv.ParseInt(process_end, 10, 32)
				length := endStart - posStart + 1
				proCode := *strCode
				proCode = proCode[posStart : posStart+length] // posStart从零还是从1开始？

				if errstrat != nil || errend != nil {
					Normal = -1
					Error = errors.New("Error:PosStart or EndStart Error in LongCode\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
					return Normal, Error
				}

				if strings.ToLower(para.ParaRangeList[k].ParaRangeSpecification) != "" {
					*strResult = para.ParaRangeList[k].ParaRangeSpecification
					*strCode = *strCode
					*strResultValue = proCode
					Normal = 0
					return Normal, Error
				} else {
					*strResult = proCode
					*strCode = *strCode
					*strResultValue = proCode
					Normal = 0
					return Normal, Error
				}
				*strResultValue = proCode
			}
		}
	} else if strings.ToLower(process_type) == "result" {
		Type := para.Type
		if strings.ToLower(Type) == "int" {
			*strCode = *strCode
			for k := 0; k < len(para.ParaRangeList); k++ {
				para.ParaRangeList[k].Alarm_max_equal = strings.ToLower(para.ParaRangeList[k].Alarm_max_equal)
				para.ParaRangeList[k].Alarm_min_equal = strings.ToLower(para.ParaRangeList[k].Alarm_min_equal)
				sup, errsup := strconv.ParseInt(para.ParaRangeList[k].Alarm_max, 10, 64)
				inf, errinf := strconv.ParseInt(para.ParaRangeList[k].Alarm_min, 10, 64)
				value, errval := strconv.ParseInt(*strCode, 16, 64)

				if errsup != nil || errinf != nil || errval != nil {
					Normal = -1
					Error = errors.New("Error:Sup, Inf or Val Error in Int\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
					return Normal, Error
				}
				if para.ParaRangeList[k].Alarm_max_equal == "true" && para.ParaRangeList[k].Alarm_min_equal == "true" {
					if value >= inf && value <= sup {
						if para.ParaRangeList[k].ParaRangeSpecification != "" {
							*strResult = para.ParaRangeList[k].ParaRangeSpecification
							*strResultValue = strconv.FormatInt(value, 10)
						} else {
							*strResultValue = strconv.FormatInt(value, 10)
						}
						Normal = 0
						return Normal, Error
					}
				} else if para.ParaRangeList[k].Alarm_max_equal == "false" && para.ParaRangeList[k].Alarm_min_equal == "true" {
					if value >= inf && value < sup {
						if para.ParaRangeList[k].ParaRangeSpecification != "" {
							*strResult = para.ParaRangeList[k].ParaRangeSpecification
							*strResultValue = strconv.FormatInt(value, 10)
						} else {
							*strResultValue = strconv.FormatInt(value, 10)
						}
						Normal = 0
						return Normal, Error
					}
				} else if para.ParaRangeList[k].Alarm_max_equal == "true" && para.ParaRangeList[k].Alarm_min_equal == "false" {
					if value > inf && value <= sup {
						if para.ParaRangeList[k].ParaRangeSpecification != "" {
							*strResult = para.ParaRangeList[k].ParaRangeSpecification
							*strResultValue = strconv.FormatInt(value, 10)
						} else {
							*strResultValue = strconv.FormatInt(value, 10)
						}
						Normal = 0
						return Normal, Error
					}
				} else if para.ParaRangeList[k].Alarm_max_equal == "false" && para.ParaRangeList[k].Alarm_min_equal == "false" {
					if value > inf && value < sup {
						if para.ParaRangeList[k].ParaRangeSpecification != "" {
							*strResult = para.ParaRangeList[k].ParaRangeSpecification
							*strResultValue = strconv.FormatInt(value, 10)
						} else {
							*strResultValue = strconv.FormatInt(value, 10)
						}
						Normal = 0
						return Normal, Error
					}
				} else {
					*strResultValue = strconv.FormatInt(value, 10)
					Normal = -1
					Error = errors.New("Fatal:Limitation Break\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
					return Normal, Error
				}
				*strResultValue = strconv.FormatInt(value, 10)
			}
		} else if strings.ToLower(Type) == "float" {
			*strCode = *strCode
			l := 0
			for k := 0; k < len(para.ParaRangeList); k++ {
				value, errval := strconv.ParseFloat(*strResult, 64)
				if errval != nil {
					*strResultValue = *strResult
					Normal = -1
					Error = errors.New("Error:Value Error in Float\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
					return Normal, Error
				}

				if para.ParaRangeList[k].Alarm_max == "limitation" && para.ParaRangeList[k].Alarm_min != "limitation" {
					inf, errinf := strconv.ParseFloat(para.ParaRangeList[k].Alarm_min, 64)

					if errinf != nil {
						Normal = -1
						Error = errors.New("Error:Inf Error in Float\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
						return Normal, Error
					}

					if strings.ToLower(para.ParaRangeList[k].Alarm_min_equal) == "true" {
						if value >= inf {
							*strResultValue = *strResult
							*strResult = para.ParaRangeList[k].ParaRangeSpecification
							if *strResult == "" {
								*strResult = *strResultValue
							}
							Normal = 0
							return Normal, Error
						}
						l++
					} else if strings.ToLower(para.ParaRangeList[k].Alarm_min_equal) == "false" {
						if value > inf {
							*strResultValue = *strResult
							*strResult = para.ParaRangeList[k].ParaRangeSpecification
							if *strResult == "" {
								*strResult = *strResultValue
							}
							Normal = 0
							return Normal, Error
						}
						l++
					}
				} else if para.ParaRangeList[k].Alarm_max != "limitation" && para.ParaRangeList[k].Alarm_min == "limitation" {
					sup, errsup := strconv.ParseFloat(para.ParaRangeList[k].Alarm_max, 64)

					if errsup != nil {
						Normal = -1
						Error = errors.New("Error:Sup Error in Float\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
						return Normal, Error
					}

					if strings.ToLower(para.ParaRangeList[k].Alarm_min_equal) == "true" {
						if value <= sup {
							*strResultValue = *strResult
							*strResult = para.ParaRangeList[k].ParaRangeSpecification
							if *strResult == "" {
								*strResult = *strResultValue
							}
							Normal = 0
							return Normal, Error
						}
						l++
					} else if strings.ToLower(para.ParaRangeList[k].Alarm_min_equal) == "false" {
						if value < sup {
							*strResultValue = *strResult
							*strResult = para.ParaRangeList[k].ParaRangeSpecification
							if *strResult == "" {
								*strResult = *strResultValue
							}
							Normal = 0
							return Normal, Error
						}
						l++
					}
				} else if strings.ToLower(para.ParaRangeList[k].Alarm_max) != "limitation" && strings.ToLower(para.ParaRangeList[k].Alarm_min) != "limitation" {
					sup, errsup := strconv.ParseFloat(para.ParaRangeList[k].Alarm_max, 32)
					inf, errinf := strconv.ParseFloat(para.ParaRangeList[k].Alarm_min, 32)

					if errsup != nil || errinf != nil {
						Normal = -1
						Error = errors.New("Error:Sup or Inf Error in Float\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
						return Normal, Error
					}
					para.ParaRangeList[k].Alarm_max_equal = strings.ToLower(para.ParaRangeList[k].Alarm_max_equal)
					para.ParaRangeList[k].Alarm_min_equal = strings.ToLower(para.ParaRangeList[k].Alarm_min_equal)

					if para.ParaRangeList[k].Alarm_max_equal == "true" && para.ParaRangeList[k].Alarm_min_equal == "true" {
						if value >= inf && value <= sup {
							*strResultValue = *strResult
							*strResult = para.ParaRangeList[k].ParaRangeSpecification
							if *strResult == "" {
								*strResult = *strResultValue
							}
							Normal = 0
							return Normal, Error
						}
						l++
					} else if para.ParaRangeList[k].Alarm_max_equal == "false" && para.ParaRangeList[k].Alarm_min_equal == "true" {
						if value >= inf && value < sup {
							*strResultValue = *strResult
							*strResult = para.ParaRangeList[k].ParaRangeSpecification
							if *strResult == "" {
								*strResult = *strResultValue
							}
							Normal = 0
							return Normal, Error
						}
						l++
					} else if para.ParaRangeList[k].Alarm_max_equal == "true" && para.ParaRangeList[k].Alarm_min_equal == "false" {
						if value > inf && value <= sup {
							*strResultValue = *strResult
							*strResult = para.ParaRangeList[k].ParaRangeSpecification
							if *strResult == "" {
								*strResult = *strResultValue
							}
							Normal = 0
							return Normal, Error
						}
						l++
					} else if para.ParaRangeList[k].Alarm_max_equal == "false" && para.ParaRangeList[k].Alarm_min_equal == "false" {
						if value > inf && value < sup {
							*strResultValue = *strResult
							*strResult = para.ParaRangeList[k].ParaRangeSpecification
							if *strResult == "" {
								*strResult = *strResultValue
							}
							Normal = 0
							return Normal, Error

						}
						l++
					} else {
						*strResultValue = strconv.FormatFloat(value, 'f', -1, 64)
						Normal = -1
						Error = errors.New("Fatal:Limitation Break\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
						return Normal, Error
					}
				}
			}
			if l > 0 {
				*strResultValue = *strResult
				Normal = -1
				Error = errors.New("Fatal:Limitation Break\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
				return Normal, Error
			} else {
				*strResultValue = *strResult
				Normal = 0
				return Normal, Error
			}
		}
	}
	Normal = -1
	Error = errors.New("Fatal:Limitation Break\nPara Name:" + para.Name + "\nParaValue:" + *strResultValue + "\n")
	return Normal, Error
}

// return 0: not changed; 1: new regist; -1: unregist
func regist_view_chan(frame *model.FrameDict, info *model.View_page_regist_info, para_view_map map[*websocket.NSConn]map[int]string) int {
	if info.Action == 1 {
		// regist only required parameters for view
		para_view_map[info.Conn] = make(map[int]string)

		var view_dict *model.ViewDict = nil
		for _, v := range *(info.View_dict) {
			if (*v).View_type.PayloadName == frame.Frame_type.PayloadName && (*v).View_type.SubAddressName == frame.Frame_type.SubAddressName {
				view_dict = v
				break
			}
		}
		if view_dict == nil {
			log.Printf("Failed register New channel no %s %s %s, now %d user\n", frame.Frame_type.MissionID, info.Conn, info.File,
				len(frame.Frame_type.UserChanMap))
			return 0
		}
		for index, item := range frame.ParaList {
			for _, p := range (*view_dict).ParaList {
				if p.ParaKey == item.ParaKey {
					(para_view_map[info.Conn])[index] = p.Index
					// log.Printf("bound %s\n", p.Index)
				}
			}
		}

		// bind from conn to view_chan
		frame.Frame_type.UserChanMap[info.Conn] = info.View_chan
		log.Printf("New channel no %s %s %s, now %d user\n", frame.Frame_type.MissionID, info.Conn, info.File,
			len(frame.Frame_type.UserChanMap))
		return 1
	} else {
		// unregist para for view, sub_map will be GC autoly
		delete(para_view_map, info.Conn)

		// unbind conn with view_chan
		delete(frame.Frame_type.UserChanMap, info.Conn)
		// 关闭chan
		close(info.View_chan)
		log.Printf("Delete channel no %s %s %s, now %d user\n", frame.Frame_type.MissionID, info.Conn, info.File,
			len(frame.Frame_type.UserChanMap))
		return -1
	}
}

func init_cases(
	chan_net_frame chan string,
	ticker *time.Ticker,
	chan_view_reg chan *model.View_page_regist_info) (cases []reflect.SelectCase) {

	// chan view register
	selectcase := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(chan_view_reg),
	}
	cases = append(cases, selectcase)

	// 定时器
	selectcase = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ticker.C),
	}
	cases = append(cases, selectcase)

	// chan frame register
	selectcase = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(chan_net_frame),
	}
	cases = append(cases, selectcase)

	return
}

func send_to_view(
	cases *[]reflect.SelectCase,
	user_chan_view_map map[*websocket.NSConn]chan string,
	send_value_map map[chan string]interface{}) {

	// 每个消费者，发送一次后必须删除
	for _, item := range user_chan_view_map {
		send_value := send_value_map[item]
		// 空字符串不发送
		if send_value.(string) == "" {
			continue
		}
		selectcase := reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(item),
			Send: reflect.ValueOf(send_value),
		}
		*cases = append(*cases, selectcase)
	}
}

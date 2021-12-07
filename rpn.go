//将一个普通的中缀表达式转换为逆波兰表达式的一般算法是：
//首先需要分配2个栈，一个作为临时存储运算符的栈S1（含一个结束符号），一个作为存放结果（逆波兰式）的栈S2（空栈），
// S1栈可先放入优先级最低的运算符#，注意，中缀式应以此最低优先级的运算符结束。可指定其他字符，不一定非#不可。从中缀式的左端开始取字符，逐序进行如下步骤：
//（1）若取出的字符是操作数，则分析出完整的运算数，该操作数直接送入S2栈。
//（2）若取出的字符是运算符，则将该运算符与S1栈栈顶元素比较，
//    如果该运算符(不包括括号运算符)优先级高于S1栈栈顶运算符（包括左括号）优先级，则将该运算符进S1栈，
//    否则，将S1栈的栈顶运算符弹出，送入S2栈中，直至S1栈栈顶运算符（包括左括号）低于（不包括等于）该运算符优先级时停止弹出运算符，最后将该运算符送入S1栈。
//（3）若取出的字符是“（”，则直接送入S1栈顶。
//（4）若取出的字符是“）”，则将距离S1栈栈顶最近的“（”之间的运算符，逐个出栈，依次送入S2栈，此时抛弃“（”。
//（5）重复上面的1~4步，直至处理完所有的输入字符。
//（6）若取出的字符是“#”，则将S1栈内所有运算符（不包括“#”），逐个出栈，依次送入S2栈。
//完成以上步骤，S2栈便为逆波兰式输出结果。不过S2应做一下逆序处理。便可以按照逆波兰式的计算方法计算了！
package rpn

import (
	"errors"
	"strconv"
)

//运算符优先级
type Level float64

const (
	LevelAddSub Level = 10.0 // 加减法优先级
	LevelMulDiv Level = 20.0 // 乘除法优先级
)

type Operator struct {
	Symbol  string
	Execute func(args ...float64) (float64, error)
	Level   Level
}

var common = &RPN{
	op: map[string]*Operator{
		"+": {Symbol: "+", Execute: func(args ...float64) (float64, error) {
			var a, b float64
			if len(args) == 1 {
				a, b = 0, args[0]
			} else {
				a, b = args[0], args[1]
			}
			return a + b, nil
		}, Level: LevelAddSub},
		"-": {Symbol: "-", Execute: func(args ...float64) (float64, error) {
			var a, b float64
			if len(args) == 1 {
				a, b = 0, args[0]
			} else {
				a, b = args[0], args[1]
			}
			return a - b, nil
		}, Level: LevelAddSub},
		"*": {Symbol: "*", Execute: func(args ...float64) (float64, error) {
			var a, b float64
			if len(args) == 1 {
				a, b = 0, args[0]
			} else {
				a, b = args[0], args[1]
			}
			return a * b, nil
		}, Level: LevelMulDiv},
		"/": {Symbol: "/", Execute: func(args ...float64) (float64, error) {
			var a, b float64
			if len(args) == 1 {
				a, b = 0, args[0]
			} else {
				a, b = args[0], args[1]
			}
			return a / b, nil
		}, Level: LevelMulDiv},
	},
}

type RPN struct {
	op map[string]*Operator
}

func NewRPN() *RPN {
	r := new(RPN)
	*r = *common
	return r
}

// AddOP 自定义运算符
// 运算符禁止包含 '(', ')', '.', 'e', 'E' 及数字
func (p *RPN) AddOP(symbol string, level Level, execute func(a ...float64) (float64, error)) *RPN {
	for _, v := range []byte(symbol) {
		if (v >= '0' && v <= '9') || v == '(' || v == ')' || v == '.' || v == 'e' || v == 'E' {
			panic("symbol cannot contain '(', ')', '.', 'E', 'e' or number")
		}
	}
	p.op[symbol] = &Operator{
		Symbol:  symbol,
		Execute: execute,
		Level:   level,
	}
	return p
}

//Parse 将中缀表达式解析为后缀表达式（逆波兰式）
func (p *RPN) Parse(exp string) ([]string, error) {
	var num string
	var op string
	var s1, s2 []string
	for _, v := range []byte(exp) {
		switch {
		case v == ' ':
			if num != "" {
				s2 = append(s2, num)
				num = ""
			}
			if op != "" {
				return nil, errors.New("invalid expression")
			}
			continue
		case v == '(':
			if num != "" {
				s2 = append(s2, num)
				num = ""
			}
			if op != "" {
				return nil, errors.New("invalid expression")
			}
			s1 = append(s1, "(")
		case v == ')':
			if num != "" {
				s2 = append(s2, num)
				num = ""
			}
			if op != "" {
				return nil, errors.New("invalid expression")
			}
			for i := len(s1) - 1; i >= 0; i-- {
				if s1[i] == "(" {
					s1 = s1[:i]
					break
				}
				s2 = append(s2, s1[i])
			}
		case p.op[string([]byte{v})] != nil:
			if num != "" {
				s2 = append(s2, num)
				num = ""
			}
			if op != "" {
				return nil, errors.New("invalid expression")
			}
			o1 := p.op[string([]byte{v})]
			if len(s1) == 0 || s1[len(s1)-1] == "(" {
				s1 = append(s1, string([]byte{v}))
				continue
			}

			o2 := p.op[s1[len(s1)-1]]
			if o1.Level > o2.Level {
				s1 = append(s1, string([]byte{v}))
			} else {
				for i := len(s1) - 1; i >= 0; i-- {
					o := p.op[s1[i]]
					if o.Level < o1.Level {
						s1 = s1[:i]
						break
					}
					s2 = append(s2, s1[i])
					if i == 0 {
						s1 = []string{}
					}
				}
				s1 = append(s1, string([]byte{v}))
			}
		case v >= '0' && v <= '9' || v == '.':
			if op != "" {
				return nil, errors.New("invalid expression")
			}
			num = string(append([]byte(num), v))
		default:
			if num != "" {
				s2 = append(s2, num)
				num = ""
			}
			op = string(append([]byte(op), v))
			if p.op[op] != nil {
				o1 := p.op[op]
				if len(s1) == 0 || s1[len(s1)-1] == "(" {
					s1 = append(s1, op)
					op = ""
					continue
				}
				o2 := p.op[s1[len(s1)-1]]
				if o1.Level > o2.Level {
					s1 = append(s1, op)
				} else {
					for i := len(s1) - 1; i >= 0; i-- {
						o := p.op[s1[i]]
						if o.Level < o1.Level {
							s1 = s1[:i]
							break
						}
						s2 = append(s2, s1[i])
						if i == 0 {
							s1 = []string{}
						}
					}
					s1 = append(s1, op)
				}
				op = ""
			}
		}
	}
	if num != "" {
		s2 = append(s2, num)
		num = ""
	}
	for i := len(s1) - 1; i >= 0; i-- {
		s2 = append(s2, s1[i])
	}
	return s2, nil
}

//WARN 仅支持双目运算符
//使用单目运算（如1+-2，1+--2）可能会出现预料外的计算结果
func (p *RPN) Calculate(s string) (float64, error) {
	exp, err := p.Parse(s)
	if err != nil {
		return 0, err
	}
	var s1 []float64
	for _, v := range exp {
		if op := p.op[v]; op != nil {
			var er error
			if len(s1) == 0 {
				return 0, errors.New("invalid expression")
			} else if len(s1) == 1 { //单目运算
				s1[0], er = op.Execute(s1[0])
			} else if len(s1) > 1 { //双目运算
				a := s1[len(s1)-2]
				b := s1[len(s1)-1]
				s1[len(s1)-2], er = op.Execute(a, b)
				s1 = s1[:len(s1)-1]
			}
			if er != nil {
				return 0, er
			}
		} else {
			n, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return 0, nil
			}
			s1 = append(s1, n)
		}
	}
	if len(s1) != 1 {
		return 0, errors.New("invalid expression")
	} else {
		return s1[0], nil
	}
}

func Calculate(s string) (float64, error) {
	return common.Calculate(s)
}

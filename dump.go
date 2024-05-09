package dump

import (
	"context"
	"fmt"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/plugin/pkg/replacer"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// Dump implement the plugin interface.
type Dump struct {
	Next plugin.Handler
}

func init() {
	fmt.Println("init")
	plugin.Register("dump", setup)
}

func setup(c *caddy.Controller) error {
	fmt.Println("setup(): start")
	for c.Next() {
		fmt.Println("setup(): c.Next()")
		if c.NextArg() {
			fmt.Println("setup(): c.NextArg()")
			return plugin.Error("dump", c.ArgErr())
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		fmt.Println("setup(): AddPlugin()")
		return Dump{Next: next}
	})
	fmt.Println("setup(): ok")
	return nil
}

var rlog = clog.NewWithPlugin("example")

const format = `{remote} ` + replacer.EmptyValue + ` {>id} {type} {class} {name} {proto} {port}`

//var output io.Writer = os.Stdout

// ServeDNS implements the plugin.Handler interface.
func (d Dump) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	rlog.Debug("ServeDNS(): Received response")
	fmt.Println("ServeDNS(): Received response")
	state := request.Request{W: w, Req: r}
	rep := replacer.New()
	trw := dnstest.NewRecorder(w)

	//fmt.Fprintln(output, rep.Replace(ctx, state, trw, format))
	fmt.Println("ServeDNS(): will log")
	fmt.Println(rep.Replace(ctx, state, trw, format))
	return plugin.NextOrFailure(d.Name(), d.Next, ctx, w, r)
}

// Name implements the Handler interface.
func (d Dump) Name() string {
	fmt.Println("Name()")
	return "dump"
}
func (d Dump) Ready() bool {
	fmt.Println("Ready()")
	return true
}

// ResponsePrinter wrap a dns.ResponseWriter and will write example to standard output when WriteMsg is called.
type ResponsePrinter struct {
	dns.ResponseWriter
}

// NewResponsePrinter returns ResponseWriter.
func NewResponsePrinter(w dns.ResponseWriter) *ResponsePrinter {
	return &ResponsePrinter{ResponseWriter: w}
}

// WriteMsg calls the underlying ResponseWriter's WriteMsg method and prints "example" to standard output.
func (r *ResponsePrinter) WriteMsg(res *dns.Msg) error {
	rlog.Info("dump")
	fmt.Println("WriteMsg(): dump")
	return r.ResponseWriter.WriteMsg(res)
}

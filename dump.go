package dump

import (
	"context"
	"log"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

// Dump implement the plugin interface.
type Dump struct {
	Next plugin.Handler
}

func init() {
	log.Printf("init")
	plugin.Register("dump", setup)
}

func setup(c *caddy.Controller) error {
	log.Printf("setup(): start")
	for c.Next() {
		log.Printf("setup(): c.Next()")
		if c.NextArg() {
			log.Printf("setup(): c.NextArg()")
			return plugin.Error("dump", c.ArgErr())
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		log.Printf("setup(): AddPlugin()")
		return &Dump{Next: next}
	})
	log.Printf("setup(): ok")
	return nil
}

//var rlog = clog.NewWithPlugin("example")

//const format = `{remote} ` + replacer.EmptyValue + ` {>id} {type} {class} {name} {proto} {port}`

//var output io.Writer = os.Stdout

// ServeDNS implements the plugin.Handler interface.
func (d *Dump) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	//rlog.Debug("ServeDNS(): Received response")
	log.Printf("ServeDNS(): Received response")
	//state := request.Request{W: w, Req: r}
	//rep := replacer.New()
	//trw := dnstest.NewRecorder(w)

	//fmt.Fprintln(output, rep.Replace(ctx, state, trw, format))
	log.Printf("ServeDNS(): will log")
	//log.Printf(rep.Replace(ctx, state, trw, format))
	return plugin.NextOrFailure(d.Name(), d.Next, ctx, w, r)
}

// Name implements the Handler interface.
func (d *Dump) Name() string {
	log.Printf("Name()")
	return "dump"
}
func (d *Dump) Ready() bool {
	log.Printf("Ready()")
	return true
}

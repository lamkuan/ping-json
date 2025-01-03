package ping

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
)

func Ping(args map[string]interface{}) (result string, err error) {
	ip, _ := args["ip"].(string)

	count, _ := args["count"].(int)

	timeout, _ := args["timeout"].(time.Duration)

	pinger, err := ping.NewPinger(ip)
	if err != nil {
		panic(err)
	}

	pinger.Count = count
	pinger.Timeout = time.Duration(timeout)

	pinger.OnRecv = func(pkt *ping.Packet) {
		result += fmt.Sprintf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
		result += fmt.Sprintf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		result += fmt.Sprintf("\n--- %s ping statistics ---\n", stats.Addr)
		result += fmt.Sprintf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		result += fmt.Sprintf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	result += fmt.Sprintf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())

	err = pinger.Run()
	if err != nil {
		panic(err)
	}

	return
}

package opnsense

type protocolStatisticsResponse struct {
	Statistics struct {
		TCP struct {
			SentPackets                           int `json:"sent-packets"`
			SentDataPackets                       int `json:"sent-data-packets"`
			SentDataBytes                         int `json:"sent-data-bytes"`
			SentRetransmittedPackets              int `json:"sent-retransmitted-packets"`
			SentRetransmittedBytes                int `json:"sent-retransmitted-bytes"`
			SentUnnecessaryRetransmittedPackets   int `json:"sent-unnecessary-retransmitted-packets"`
			SentResendsByMtuDiscovery             int `json:"sent-resends-by-mtu-discovery"`
			SentAckOnlyPackets                    int `json:"sent-ack-only-packets"`
			SentPacketsDelayed                    int `json:"sent-packets-delayed"`
			SentUrgOnlyPackets                    int `json:"sent-urg-only-packets"`
			SentWindowProbePackets                int `json:"sent-window-probe-packets"`
			SentWindowUpdatePackets               int `json:"sent-window-update-packets"`
			SentControlPackets                    int `json:"sent-control-packets"`
			ReceivedPackets                       int `json:"received-packets"`
			ReceivedAckPackets                    int `json:"received-ack-packets"`
			ReceivedAckBytes                      int `json:"received-ack-bytes"`
			ReceivedDuplicateAcks                 int `json:"received-duplicate-acks"`
			ReceivedUDPTunneledPkts               int `json:"received-udp-tunneled-pkts"`
			ReceivedBadUDPTunneledPkts            int `json:"received-bad-udp-tunneled-pkts"`
			ReceivedAcksForUnsentData             int `json:"received-acks-for-unsent-data"`
			ReceivedInSequencePackets             int `json:"received-in-sequence-packets"`
			ReceivedInSequenceBytes               int `json:"received-in-sequence-bytes"`
			ReceivedCompletelyDuplicatePackets    int `json:"received-completely-duplicate-packets"`
			ReceivedCompletelyDuplicateBytes      int `json:"received-completely-duplicate-bytes"`
			ReceivedOldDuplicatePackets           int `json:"received-old-duplicate-packets"`
			ReceivedSomeDuplicatePackets          int `json:"received-some-duplicate-packets"`
			ReceivedSomeDuplicateBytes            int `json:"received-some-duplicate-bytes"`
			ReceivedOutOfOrder                    int `json:"received-out-of-order"`
			ReceivedOutOfOrderBytes               int `json:"received-out-of-order-bytes"`
			ReceivedAfterWindowPackets            int `json:"received-after-window-packets"`
			ReceivedAfterWindowBytes              int `json:"received-after-window-bytes"`
			ReceivedWindowProbes                  int `json:"received-window-probes"`
			ReceiveWindowUpdatePackets            int `json:"receive-window-update-packets"`
			ReceivedAfterClosePackets             int `json:"received-after-close-packets"`
			DiscardBadChecksum                    int `json:"discard-bad-checksum"`
			DiscardBadHeaderOffset                int `json:"discard-bad-header-offset"`
			DiscardTooShort                       int `json:"discard-too-short"`
			DiscardReassemblyQueueFull            int `json:"discard-reassembly-queue-full"`
			ConnectionRequests                    int `json:"connection-requests"`
			ConnectionsAccepts                    int `json:"connections-accepts"`
			BadConnectionAttempts                 int `json:"bad-connection-attempts"`
			ListenQueueOverflows                  int `json:"listen-queue-overflows"`
			IgnoredInWindowResets                 int `json:"ignored-in-window-resets"`
			ConnectionsEstablished                int `json:"connections-established"`
			ConnectionsHostcacheRtt               int `json:"connections-hostcache-rtt"`
			ConnectionsHostcacheRttvar            int `json:"connections-hostcache-rttvar"`
			ConnectionsHostcacheSsthresh          int `json:"connections-hostcache-ssthresh"`
			ConnectionsClosed                     int `json:"connections-closed"`
			ConnectionDrops                       int `json:"connection-drops"`
			ConnectionsUpdatedRttOnClose          int `json:"connections-updated-rtt-on-close"`
			ConnectionsUpdatedVarianceOnClose     int `json:"connections-updated-variance-on-close"`
			ConnectionsUpdatedSsthreshOnClose     int `json:"connections-updated-ssthresh-on-close"`
			EmbryonicConnectionsDropped           int `json:"embryonic-connections-dropped"`
			SegmentsUpdatedRtt                    int `json:"segments-updated-rtt"`
			SegmentUpdateAttempts                 int `json:"segment-update-attempts"`
			RetransmitTimeouts                    int `json:"retransmit-timeouts"`
			ConnectionsDroppedByRetransmitTimeout int `json:"connections-dropped-by-retransmit-timeout"`
			PersistTimeout                        int `json:"persist-timeout"`
			ConnectionsDroppedByPersistTimeout    int `json:"connections-dropped-by-persist-timeout"`
			ConnectionsDroppedByFinwait2Timeout   int `json:"connections-dropped-by-finwait2-timeout"`
			KeepaliveTimeout                      int `json:"keepalive-timeout"`
			KeepaliveProbes                       int `json:"keepalive-probes"`
			ConnectionsDroppedByKeepalives        int `json:"connections-dropped-by-keepalives"`
			AckHeaderPredictions                  int `json:"ack-header-predictions"`
			DataPacketHeaderPredictions           int `json:"data-packet-header-predictions"`
			Syncache                              struct {
				EntriesAdded   int `json:"entries-added"`
				Retransmitted  int `json:"retransmitted"`
				Duplicates     int `json:"duplicates"`
				Dropped        int `json:"dropped"`
				Completed      int `json:"completed"`
				BucketOverflow int `json:"bucket-overflow"`
				CacheOverflow  int `json:"cache-overflow"`
				Reset          int `json:"reset"`
				Stale          int `json:"stale"`
				Aborted        int `json:"aborted"`
				BadAck         int `json:"bad-ack"`
				Unreachable    int `json:"unreachable"`
				ZoneFailures   int `json:"zone-failures"`
				SentCookies    int `json:"sent-cookies"`
				ReceivdCookies int `json:"receivd-cookies"`
			} `json:"syncache"`
			Hostcache struct {
				EntriesAdded    int `json:"entries-added"`
				BufferOverflows int `json:"buffer-overflows"`
			} `json:"hostcache"`
			Sack struct {
				RecoveryEpisodes    int `json:"recovery-episodes"`
				SegmentRetransmits  int `json:"segment-retransmits"`
				ByteRetransmits     int `json:"byte-retransmits"`
				ReceivedBlocks      int `json:"received-blocks"`
				SentOptionBlocks    int `json:"sent-option-blocks"`
				ScoreboardOverflows int `json:"scoreboard-overflows"`
			} `json:"sack"`
			Ecn struct {
				CePackets            int `json:"ce-packets"`
				Ect0Packets          int `json:"ect0-packets"`
				Ect1Packets          int `json:"ect1-packets"`
				Handshakes           int `json:"handshakes"`
				CongestionReductions int `json:"congestion-reductions"`
			} `json:"ecn"`
			TCPSignature struct {
				ReceivedGoodSignature int `json:"received-good-signature"`
				ReceivedBadSignature  int `json:"received-bad-signature"`
				FailedMakeSignature   int `json:"failed-make-signature"`
				NoSignatureExpected   int `json:"no-signature-expected"`
				NoSignatureProvided   int `json:"no-signature-provided"`
			} `json:"tcp-signature"`
			Pmtud struct {
				PmtudActivated       int `json:"pmtud-activated"`
				PmtudActivatedMinMss int `json:"pmtud-activated-min-mss"`
				PmtudFailed          int `json:"pmtud-failed"`
			} `json:"pmtud"`
			Tw struct {
				TwResponds int `json:"tw_responds"`
				TwRecycles int `json:"tw_recycles"`
				TwResets   int `json:"tw_resets"`
			} `json:"tw"`
			TCPConnectionCountByState struct {
				Closed      int `json:"CLOSED"`
				Listen      int `json:"LISTEN"`
				SynSent     int `json:"SYN_SENT"`
				SynRcvd     int `json:"SYN_RCVD"`
				Established int `json:"ESTABLISHED"`
				CloseWait   int `json:"CLOSE_WAIT"`
				FinWait1    int `json:"FIN_WAIT_1"`
				Closing     int `json:"CLOSING"`
				LastAck     int `json:"LAST_ACK"`
				FinWait2    int `json:"FIN_WAIT_2"`
				TimeWait    int `json:"TIME_WAIT"`
			} `json:"TCP connection count by state"`
		} `json:"tcp"`
		UDP struct {
			ReceivedDatagrams            int `json:"received-datagrams"`
			DroppedIncompleteHeaders     int `json:"dropped-incomplete-headers"`
			DroppedBadDataLength         int `json:"dropped-bad-data-length"`
			DroppedBadChecksum           int `json:"dropped-bad-checksum"`
			DroppedNoChecksum            int `json:"dropped-no-checksum"`
			DroppedNoSocket              int `json:"dropped-no-socket"`
			DroppedBroadcastMulticast    int `json:"dropped-broadcast-multicast"`
			DroppedFullSocketBuffer      int `json:"dropped-full-socket-buffer"`
			NotForHashedPcb              int `json:"not-for-hashed-pcb"`
			DeliveredPackets             int `json:"delivered-packets"`
			OutputPackets                int `json:"output-packets"`
			MulticastSourceFilterMatches int `json:"multicast-source-filter-matches"`
		} `json:"udp"`
		IP struct {
			ReceivedPackets               int `json:"received-packets"`
			DroppedBadChecksum            int `json:"dropped-bad-checksum"`
			DroppedBelowMinimumSize       int `json:"dropped-below-minimum-size"`
			DroppedShortPackets           int `json:"dropped-short-packets"`
			DroppedTooLong                int `json:"dropped-too-long"`
			DroppedShortHeaderLength      int `json:"dropped-short-header-length"`
			DroppedShortData              int `json:"dropped-short-data"`
			DroppedBadOptions             int `json:"dropped-bad-options"`
			DroppedBadVersion             int `json:"dropped-bad-version"`
			ReceivedFragments             int `json:"received-fragments"`
			DroppedFragments              int `json:"dropped-fragments"`
			DroppedFragmentsAfterTimeout  int `json:"dropped-fragments-after-timeout"`
			ReassembledPackets            int `json:"reassembled-packets"`
			ReceivedLocalPackets          int `json:"received-local-packets"`
			DroppedUnknownProtocol        int `json:"dropped-unknown-protocol"`
			ForwardedPackets              int `json:"forwarded-packets"`
			FastForwardedPackets          int `json:"fast-forwarded-packets"`
			PacketsCannotForward          int `json:"packets-cannot-forward"`
			ReceivedUnknownMulticastGroup int `json:"received-unknown-multicast-group"`
			RedirectsSent                 int `json:"redirects-sent"`
			SentPackets                   int `json:"sent-packets"`
			SendPacketsFabricatedHeader   int `json:"send-packets-fabricated-header"`
			DiscardNoMbufs                int `json:"discard-no-mbufs"`
			DiscardNoRoute                int `json:"discard-no-route"`
			SentFragments                 int `json:"sent-fragments"`
			FragmentsCreated              int `json:"fragments-created"`
			DiscardCannotFragment         int `json:"discard-cannot-fragment"`
			DiscardTunnelNoGif            int `json:"discard-tunnel-no-gif"`
			DiscardBadAddress             int `json:"discard-bad-address"`
		} `json:"ip"`
		Icmp struct {
			IcmpCalls            int `json:"icmp-calls"`
			ErrorsNotFromMessage int `json:"errors-not-from-message"`
			OutputHistogram      []struct {
				Name  string `json:"name"`
				Count int    `json:"count"`
			} `json:"output-histogram"`
			DroppedBadCode            int `json:"dropped-bad-code"`
			DroppedTooShort           int `json:"dropped-too-short"`
			DroppedBadChecksum        int `json:"dropped-bad-checksum"`
			DroppedBadLength          int `json:"dropped-bad-length"`
			DroppedMulticastEcho      int `json:"dropped-multicast-echo"`
			DroppedMulticastTimestamp int `json:"dropped-multicast-timestamp"`
			InputHistogram            []struct {
				Name  string `json:"name"`
				Count int    `json:"count"`
			} `json:"input-histogram"`
			SentPackets                 int    `json:"sent-packets"`
			DiscardInvalidReturnAddress int    `json:"discard-invalid-return-address"`
			DiscardNoRoute              int    `json:"discard-no-route"`
			IcmpAddressResponses        string `json:"icmp-address-responses"`
		} `json:"icmp"`
		Carp struct {
			ReceivedInetPackets      int `json:"received-inet-packets"`
			ReceivedInet6Packets     int `json:"received-inet6-packets"`
			DroppedWrongTTL          int `json:"dropped-wrong-ttl"`
			DroppedShortHeader       int `json:"dropped-short-header"`
			DroppedBadChecksum       int `json:"dropped-bad-checksum"`
			DroppedBadVersion        int `json:"dropped-bad-version"`
			DroppedShortPacket       int `json:"dropped-short-packet"`
			DroppedBadAuthentication int `json:"dropped-bad-authentication"`
			DroppedBadVhid           int `json:"dropped-bad-vhid"`
			DroppedBadAddressList    int `json:"dropped-bad-address-list"`
			SentInetPackets          int `json:"sent-inet-packets"`
			SentInet6Packets         int `json:"sent-inet6-packets"`
			SendFailedMemoryError    int `json:"send-failed-memory-error"`
		} `json:"carp"`
		Pfsync struct {
			ReceivedInetPackets  int `json:"received-inet-packets"`
			ReceivedInet6Packets int `json:"received-inet6-packets"`
			InputHistogram       []struct {
				Name  string `json:"name"`
				Count int    `json:"count"`
			} `json:"input-histogram"`
			DroppedBadInterface int `json:"dropped-bad-interface"`
			DroppedBadTTL       int `json:"dropped-bad-ttl"`
			DroppedShortHeader  int `json:"dropped-short-header"`
			DroppedBadVersion   int `json:"dropped-bad-version"`
			DroppedBadAuth      int `json:"dropped-bad-auth"`
			DroppedBadAction    int `json:"dropped-bad-action"`
			DroppedShort        int `json:"dropped-short"`
			DroppedBadValues    int `json:"dropped-bad-values"`
			DroppedStaleState   int `json:"dropped-stale-state"`
			DroppedFailedLookup int `json:"dropped-failed-lookup"`
			SentInetPackets     int `json:"sent-inet-packets"`
			SendInet6Packets    int `json:"send-inet6-packets"`
			OutputHistogram     []struct {
				Name  string `json:"name"`
				Count int    `json:"count"`
			} `json:"output-histogram"`
			DiscardedNoMemory int `json:"discarded-no-memory"`
			SendErrors        int `json:"send-errors"`
		} `json:"pfsync"`
		Arp struct {
			SentRequests            int `json:"sent-requests"`
			SentFailures            int `json:"sent-failures"`
			SentReplies             int `json:"sent-replies"`
			ReceivedRequests        int `json:"received-requests"`
			ReceivedReplies         int `json:"received-replies"`
			ReceivedPackets         int `json:"received-packets"`
			DroppedNoEntry          int `json:"dropped-no-entry"`
			EntriesTimeout          int `json:"entries-timeout"`
			DroppedDuplicateAddress int `json:"dropped-duplicate-address"`
		} `json:"arp"`
	} `json:"statistics"`
}

type ProtocolStatistics struct {
	TCPSentPackets            int
	TCPReceivedPackets        int
	ARPSentRequests           int
	ARPReceivedRequests       int
	TCPConnectionCountByState map[string]int
}

func (c *Client) FetchProtocolStatistics() (ProtocolStatistics, *APICallError) {
	var (
		resp protocolStatisticsResponse
	)
	url, ok := c.endpoints["protocolStatistics"]
	if !ok {
		return ProtocolStatistics{}, &APICallError{
			Endpoint:   "protocolStatistics",
			StatusCode: 404,
			Message:    "endpoint not found in client endpoints",
		}
	}
	if err := c.do("GET", url, nil, &resp); err != nil {
		return ProtocolStatistics{}, err
	}

	out := ProtocolStatistics{
		TCPSentPackets:      resp.Statistics.TCP.SentPackets,
		TCPReceivedPackets:  resp.Statistics.TCP.ReceivedPackets,
		ARPSentRequests:     resp.Statistics.Arp.SentRequests,
		ARPReceivedRequests: resp.Statistics.Arp.ReceivedRequests,
		TCPConnectionCountByState: map[string]int{
			"CLOSED":      resp.Statistics.TCP.TCPConnectionCountByState.Closed,
			"LISTEN":      resp.Statistics.TCP.TCPConnectionCountByState.Listen,
			"SYN_SENT":    resp.Statistics.TCP.TCPConnectionCountByState.SynSent,
			"SYN_RCVD":    resp.Statistics.TCP.TCPConnectionCountByState.SynRcvd,
			"ESTABLISHED": resp.Statistics.TCP.TCPConnectionCountByState.Established,
			"CLOSE_WAIT":  resp.Statistics.TCP.TCPConnectionCountByState.CloseWait,
			"FIN_WAIT_1":  resp.Statistics.TCP.TCPConnectionCountByState.FinWait1,
			"CLOSING":     resp.Statistics.TCP.TCPConnectionCountByState.Closing,
			"LAST_ACK":    resp.Statistics.TCP.TCPConnectionCountByState.LastAck,
			"FIN_WAIT_2":  resp.Statistics.TCP.TCPConnectionCountByState.FinWait2,
			"TIME_WAIT":   resp.Statistics.TCP.TCPConnectionCountByState.TimeWait,
		},
	}

	return out, nil
}

/*
   Copyright 2014 Krishna Raman <kraman@gmail.com>
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
       http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

#ifndef _NETFILTER_H
#define _NETFILTER_H

#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <unistd.h>
#include <netinet/in.h>
#include <linux/types.h>
#include <linux/socket.h>
#include <linux/netfilter.h>
#include <errno.h>
#include <libnetfilter_queue/libnetfilter_queue.h>

extern uint processPacket(int id, unsigned char* data, int len, unsigned char* newData, u_int32_t idx);


// Callback for the nf handler. Follows standard message. Passes control to the Go call back
// after allocating buffer space in the C side and then sets the final verdict after the packets
// has been processed.
static int nf_callback(struct nfq_q_handle *qh, struct nfgenmsg *nfmsg, struct nfq_data *nfa, void *cb_func){
    uint32_t id = -1;
    struct nfqnl_msg_packet_hdr *ph = NULL;
    unsigned char *buffer , *new_data;
    int buffer_length ;

    ph = nfq_get_msg_packet_hdr(nfa);
    if (ph == NULL ) {
      // Something is going wrong. Stop processing packets
      return -1;
    }

    id = ntohl(ph->packet_id);

    buffer_length = nfq_get_payload(nfa, &buffer);

    new_data = (unsigned char *) malloc(buffer_length+1440);

    return processPacket(id, buffer, buffer_length, new_data, (uint32_t)((uintptr_t)cb_func) );

}

// CreateQueue uses the nfq API to create a queue
static inline struct nfq_q_handle* CreateQueue(struct nfq_handle *h, u_int16_t queue, u_int32_t idx)
{
  return nfq_create_queue(h, queue, &nf_callback, (void*)((uintptr_t)idx));
}

// Loop reading from the handle
static inline void Run(struct nfq_handle *h, int fd)
{
    char buf[65000] __attribute__ ((aligned));
    int rv;

    for (;;) {
      rv = recv(fd, buf, sizeof(buf), 0);
      if (rv >= 0) {
        nfq_handle_packet(h, buf, rv) ;
      } else {
        printf("Error in receiving packets from NFQUEUE %d\n", queue )
      }
    }

}

// SetVerdict use the nfq API to set the verdict
static inline int SetVerdict(struct nfq_q_handle *qh, int id, int verdict , int buffer_length , unsigned char *buffer ) {
    return nfq_set_verdict(qh, id , verdict  , buffer_length, buffer );
}

#endif

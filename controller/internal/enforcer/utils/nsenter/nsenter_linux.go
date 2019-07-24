// +build linux

package nsenter

/*

#cgo CFLAGS: -Wall
#include<stdio.h>
#include<sys/stat.h>
#include<sys/types.h>
#include<errno.h>
#include<string.h>
extern int errno;
extern void nsexec();
extern void droppriveleges();
extern void setupiptables();
void __attribute__((constructor)) init(void) {

	nsexec();
        setupiptables();
        droppriveleges();
}
*/
import "C"

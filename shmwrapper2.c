#include <stdio.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>

#define SOCKET_FD 1
#define SEND_FD 0

int main(int argc, char** argv) {
    /* Send FD 0 over the open socket in FD 1 */
    struct msghdr msg = {0};
    struct cmsghdr *cmsg;
    char buf[CMSG_SPACE(sizeof(int))];
    int *fdptr;
    msg.msg_control = buf;
    msg.msg_controllen = sizeof(buf);
    cmsg = CMSG_FIRSTHDR(&msg);
    cmsg->cmsg_level = SOL_SOCKET;
    cmsg->cmsg_type = SCM_RIGHTS;
    cmsg->cmsg_len = CMSG_LEN(sizeof(int));
    fdptr = (int *) CMSG_DATA(cmsg);
    fdptr[0] = SEND_FD;
    msg.msg_controllen = cmsg->cmsg_len;
    if (sendmsg(SOCKET_FD, &msg, 0) < 0)
        perror("sendmsg");
    return 0;
}

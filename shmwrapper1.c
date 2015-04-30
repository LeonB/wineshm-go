#include <windows.h>
#include <stdio.h>

int main(int argc, char** argv) {
    char *access_mode;
    DWORD access = 0;
    STARTUPINFO si;
    PROCESS_INFORMATION pi;
    HANDLE maph;

    if (argc < 4) {
        fprintf(stderr, "not enough arguments\n");
        return 1;
    }
    access_mode = argv[2];
    while (*access_mode) {
        switch(*access_mode) {
            case 'r': access |= FILE_MAP_READ; break;
            case 'w': access |= FILE_MAP_WRITE; break;
        }
        access_mode++;
    }

    maph = OpenFileMapping(access, TRUE, argv[1]);
    /* printf("maph: %p\n", maph); */
    if (maph == NULL) {
        fprintf(stderr, "failed to open mapping %s: %s\n", argv[1], strerror(GetLastError()));
        return 1;
    }
    SetStdHandle(STD_INPUT_HANDLE, maph);
    ZeroMemory(&si, sizeof(si));
    si.cb = sizeof(si);
    ZeroMemory(&pi, sizeof(pi));

    /* We can't wait for the helper to exit, since waiting for a spawned unix
     * process does not work:
     * https://bugs.winehq.org/show_bug.cgi?id=22338 */

    if (!CreateProcess(argv[3], NULL, NULL, NULL, TRUE, 0, NULL, NULL, &si, &pi)) {
        fprintf(stderr, "failed to launch second helper process: %s\n", strerror(GetLastError()));
    }
    return 0;
}

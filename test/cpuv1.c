#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main()
{
    int i = 0;
    long n = 645000000;
    while(1) {
        for(; i<n*2; i++);
        sleep(1);
    }

    return 0;
}

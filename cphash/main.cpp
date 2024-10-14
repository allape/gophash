#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include "CImg.h"

using namespace cimg_library;

extern "C" {
    double ph_hammingdistance2(uint8_t *hashA, int lenA, uint8_t *hashB, int lenB);
    uint8_t *ph_mh_imagehash(const char *filename, int *N, float alpha, float lvl);
}

int main(int argc, char **argv) {
    if (argc < 2) {
        printf("no enough input args\n");
        exit(1);
    }

    const char *img1 = argv[1];

    float alpha = 2.0;
    float level = 1.0;

    int hashlen1 = 0;

    uint8_t *hash1 = ph_mh_imagehash(img1, &hashlen1, alpha, level);

    for (int i = 0; i < hashlen1; i++) {
        // print as hex
        printf("%02x", hash1[i]);
    }

    printf("\n");

    free(hash1);
    return 0;
}


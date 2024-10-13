// ALL CODE ARE COPIED FROM https://github.com/aetilius/pHash

/*

    pHash, the open source perceptual hash library
    Copyright (C) 2009 Aetilius, Inc.
    All rights reserved.

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.

    Evan Klinger - eklinger@phash.org
    D Grant Starkweather - dstarkweather@phash.org

*/

// use PNG only, image should be preprocessed with ffmpeg
#define cimg_use_png 1
//#define cimg_use_jpeg 1
//#define cimg_use_tiff 1
#define cimg_debug 0
#define cimg_display 0

#include "CImg.h"
#include <cstdint>
#include <cmath>

using namespace std;
using namespace cimg_library;

int ph_bitcount8(uint8_t val) {
    int num = 0;
    while (val) {
        ++num;
        val &= val - 1;
    }
    return num;
}

CImg<float> *GetMHKernel(float alpha, float level) {
    int sigma = (int)4 * pow((float)alpha, (float)level);
    static CImg<float> *pkernel = NULL;
    float xpos, ypos, A;
    if (!pkernel) {
        pkernel = new CImg<float>(2 * sigma + 1, 2 * sigma + 1, 1, 1, 0);
        cimg_forXY(*pkernel, X, Y) {
            xpos = pow(alpha, -level) * (X - sigma);
            ypos = pow(alpha, -level) * (Y - sigma);
            A = xpos * xpos + ypos * ypos;
            pkernel->atXY(X, Y) = (2 - A) * exp(-A / 2);
        }
    }
    return pkernel;
}

extern "C" {
    double ph_hammingdistance2(uint8_t *hashA, int lenA, uint8_t *hashB, int lenB) {
        if (lenA != lenB) {
            return -1.0;
        }
        if ((hashA == NULL) || (hashB == NULL) || (lenA <= 0)) {
            return -1.0;
        }
        double dist = 0;
        uint8_t D = 0;
        for (int i = 0; i < lenA; i++) {
            D = hashA[i] ^ hashB[i];
            dist = dist + (double)ph_bitcount8(D);
        }
        double bits = (double)lenA * 8;
        return dist / bits;
    }

    uint8_t *ph_mh_imagehash(const char *filename, int &N, float alpha, float lvl) {
        if (filename == NULL) {
            return NULL;
        }
        uint8_t *hash = (unsigned char *)malloc(72 * sizeof(uint8_t));
        N = 72;

        CImg<uint8_t> src(filename);
        CImg<uint8_t> img;

        if (src.spectrum() == 3) {
            img = src.get_RGBtoYCbCr()
                      .channel(0)
                      .blur(1.0)
                      .resize(512, 512, 1, 1, 5)
                      .get_equalize(256);
        } else {
            img = src.channel(0)
                      .get_blur(1.0)
                      .resize(512, 512, 1, 1, 5)
                      .get_equalize(256);
        }
        src.clear();

        CImg<float> *pkernel = GetMHKernel(alpha, lvl);
        CImg<float> fresp = img.get_correlate(*pkernel);
        img.clear();
        fresp.normalize(0, 1.0);
        CImg<float> blocks(31, 31, 1, 1, 0);
        for (int rindex = 0; rindex < 31; rindex++) {
            for (int cindex = 0; cindex < 31; cindex++) {
                blocks(rindex, cindex) =
                    fresp
                        .get_crop(rindex * 16, cindex * 16, rindex * 16 + 16 - 1,
                                  cindex * 16 + 16 - 1)
                        .sum();
            }
        }
        int hash_index;
        int nb_ones = 0, nb_zeros = 0;
        int bit_index = 0;
        unsigned char hashbyte = 0;
        for (int rindex = 0; rindex < 31 - 2; rindex += 4) {
            CImg<float> subsec;
            for (int cindex = 0; cindex < 31 - 2; cindex += 4) {
                subsec = blocks.get_crop(cindex, rindex, cindex + 2, rindex + 2)
                             .unroll('x');
                float ave = subsec.mean();
                cimg_forX(subsec, I) {
                    hashbyte <<= 1;
                    if (subsec(I) > ave) {
                        hashbyte |= 0x01;
                        nb_ones++;
                    } else {
                        nb_zeros++;
                    }
                    bit_index++;
                    if ((bit_index % 8) == 0) {
                        hash_index = (int)(bit_index / 8) - 1;
                        hash[hash_index] = hashbyte;
                        hashbyte = 0x00;
                    }
                }
            }
        }

        return hash;
    }
}

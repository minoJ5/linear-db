#include <stdio.h>
#include "match.h"

int main() {
  char *needle = "am user";
  char *haystack = "app models order31.rb";
  size_t *pos;
  score_t score;
  score = match_positions(needle, haystack, pos);
  for (int i = 0; i < 8; i++) {
    printf("%ld\t", pos[i]);
  }
  printf("\n");
  printf("%f\n", score);
}
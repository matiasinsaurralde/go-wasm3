#include <stdlib.h>
#include <string.h>

char* somecall() {
  char* test = (char*) malloc(12*sizeof(char));
  strcpy(test, "testingonly");
  return test;
};


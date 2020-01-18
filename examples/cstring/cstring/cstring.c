#include <stdlib.h>
#include <string.h>

char* allocate(int n) {
  char* ptr = (char*) malloc(n*sizeof(char));
  return ptr;
}

char* somecall() {
  char* test = (char*) malloc(12*sizeof(char));
  strcpy(test, "testingonly");
  return test;
};


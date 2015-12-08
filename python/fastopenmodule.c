#include <Python.h>
#include <sys/socket.h>
#include <stdio.h>
#include <netdb.h>
#include <strings.h>
#include <stdlib.h>

static PyObject* fastopen_connect(PyObject* self, PyObject* args) {
    int socketfd;
    const char* hostname;
    const char* port;
    struct addrinfo *server;
	if (!PyArg_ParseTuple(args, "iss", &socketfd, &hostname, &port)) {
		return NULL;
	}
    
    getaddrinfo(hostname, port, 0, &server);
    
    sa_endpoints_t endpoints;
    bzero((char*)&endpoints, sizeof(endpoints));
    endpoints.sae_dstaddr = server->ai_addr;
    endpoints.sae_dstaddrlen = server->ai_addrlen;

    int rc = connectx(socketfd,
            &endpoints,
            SAE_ASSOCID_ANY,
            CONNECT_RESUME_ON_READ_WRITE | CONNECT_DATA_IDEMPOTENT,
            NULL, 0, NULL, NULL);

    freeaddrinfo(server);
    return PyLong_FromLong(rc);
}

static PyMethodDef methods[] = {
    {"connect", fastopen_connect, METH_VARARGS, "connect TCP with fastopen enabled"},
    {NULL, NULL, 0, NULL}
};
#if PY_MAJOR_VERSION >= 3
static struct PyModuleDef fastopenmodule = {
    PyModuleDef_HEAD_INIT,
    "fastopen",
    NULL,
    -1,
    methods
};
PyMODINIT_FUNC PyInit_fastopen(void)
#else
void initfastopen(void)
#endif
{
#if PY_MAJOR_VERSION >= 3
    return PyModule_Create(&fastopenmodule);
#else
    Py_InitModule("fastopen", methods);
#endif
}

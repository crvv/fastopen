#include <Python.h>
#include <sys/socket.h>
#include <stdio.h>
#include <netdb.h>
#include <strings.h>
#include <unistd.h>

static PyObject* fastopen_connect(PyObject* self, PyObject* args) {
    int socketfd;
	const char* hostname;
	int port;

	if (!PyArg_ParseTuple(args, "isi", &socketfd, &hostname, &port)) {
		return NULL;
	}

    struct hostent* server = gethostbyname(hostname);

    struct sockaddr_in server_addr;
    bzero((char*)&server_addr, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    bcopy((char*)server->h_addr_list[0], (char*)&server_addr.sin_addr.s_addr, server->h_length);
    server_addr.sin_port = htons(port);
    server_addr.sin_len = sizeof(server_addr);

    sa_endpoints_t endpoints;
    bzero((char*)&endpoints, sizeof(endpoints));
    endpoints.sae_dstaddr = (struct sockaddr*)&server_addr;
    endpoints.sae_dstaddrlen = sizeof(server_addr);

    int rc = connectx(socketfd,
            &endpoints,
            SAE_ASSOCID_ANY,
            CONNECT_RESUME_ON_READ_WRITE | CONNECT_DATA_IDEMPOTENT,
            NULL, 0, NULL, NULL);

    return PyLong_FromLong(rc);
}

static PyMethodDef methods[] = {
    {"connect", fastopen_connect, METH_VARARGS, "connect TCP with fastopen enabled"},
    {NULL, NULL, 0, NULL}
};

static struct PyModuleDef fastopenmodule = {
    PyModuleDef_HEAD_INIT,
    "fastopen",
    NULL,
    -1,
    methods
};

PyMODINIT_FUNC PyInit_fastopen(void) {
    return PyModule_Create(&fastopenmodule);
}

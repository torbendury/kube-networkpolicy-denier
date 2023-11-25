import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  insecureSkipTLSVerify: true,
  stages: [
    { duration: '15s', target: 10 },   // 10 VUs for 15s
    { duration: '15s', target: 100 },  // 100 VUs for 15s
    { duration: '30s', target: 1000 }, // 1000 VUs for 30s
    { duration: '60s', target: 10000 }, // 10000 VUs for 60s
    { duration: '15s', target: 12500 }, // 15000 VUs for 60s, errors might happen.
    { duration: '15s', target: 15000 }, // errors are expected from here on.
    { duration: '10s', target: 17500 },
    // { duration: '5s', target: 20000 },  VUs this high are not really possible locally due to broken pipes.
    // { duration: '5s', target: 25000 },
    // { duration: '5s', target: 30000 },
    // { duration: '5s', target: 35000 },
    // { duration: '5s', target: 40000 },
  ],
};

export default function () {
  const res = http.post('https://localhost:8443/validate', '{"request":{"uid":"","kind":{"group":"networking.k8s.io","version":"v1","kind":"NetworkPolicy"},"resource":{"group":"","version":"","resource":""},"operation":"CREATE","userInfo":{},"object":{"apiVersion":"networking.k8s.io/v1","kind":"NetworkPolicy","metadata":{"name":"test-policy"},"spec":{"podSelector":{"matchLabels":{"app":"test"}},"ingress":[{"from":[{"namespaceSelector":{"matchLabels":{"name":"test-ns"}}}]}]}},"oldObject":null,"options":null}}');
  check(res, { 'status was 200': (r) => r.status == 200 });
  sleep(1);
}

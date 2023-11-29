import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  insecureSkipTLSVerify: true,
  stages: [
    { duration: '10s', target: 10 },
    { duration: '30s', target: 100 },
    { duration: '60s', target: 1000 },
  ],
};

export default function () {
  const res = http.post('https://localhost:8443/validate', '{"request":{"uid":"","kind":{"group":"networking.k8s.io","version":"v1","kind":"NetworkPolicy"},"resource":{"group":"","version":"","resource":""},"operation":"CREATE","userInfo":{},"object":{"apiVersion":"networking.k8s.io/v1","kind":"NetworkPolicy","metadata":{"name":"test-policy"},"spec":{"podSelector":{"matchLabels":{"app":"test"}},"ingress":[{"from":[{"namespaceSelector":{"matchLabels":{"name":"test-ns"}}}]}]}},"oldObject":null,"options":null}}');
  check(res, { 'status was 200': (r) => r.status == 200 });
  sleep(1);
}

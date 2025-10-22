import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
  const url = 'http://localhost:8080/api/v1/data/shorten';
  const payload = JSON.stringify({ url: 'https://example.com' });
  const params = { headers: { 'Content-Type': 'application/json' } };
  http.post(url, payload, params);
  sleep(1);
}

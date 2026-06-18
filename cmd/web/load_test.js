import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 3,          
  duration: '10s',   
};

const BASE_URL = 'http://localhost:8180';

const users = [
  { email: 'sidorovich@mail.by',   password: 'password123' },   // 1
  { email: 'petrovich@mail.by',    password: 'secure456' },     // 2
  { email: 'ivanovich@mail.by',    password: 'pa$$word789' },   // 3
  { email: 'smirnova@mail.by',     password: 'catlover321' },   // 4
  { email: 'kuznecov@mail.by',     password: 'user5pass' },     // 5
  { email: 'sokolova@mail.by',     password: 'securePass6' },   // 6
  { email: 'novikov@mail.by',      password: '7novikov7' },     // 7
  { email: 'morozova@mail.by',     password: 'f80st@cold' },    // 8
  { email: 'volkov@mail.by',       password: 'w0lfPass' },      // 9
  { email: 'lebed@mail.by',        password: 'swan2023' },      // 10
  { email: 'kozlov@mail.by',       password: 'g0@tMilk' },      // 11
  { email: 'novik@mail.by',        password: 'newgirl12' },     // 12
  { email: 'solovey@mail.by',      password: 'nightingale' },   // 13
  { email: 'vasilevich@mail.by',   password: 'JulVas1988' },    // 14
  { email: 'zaycev@mail.by',       password: 'fastRunner' },    // 15
];


const transfers = [
  { to_account: 'BY20BAPB30141000000000000000', amount: 500 },
  { to_account: 'BY20BAPB30141000000000000001', amount: 800 },
  { to_account: 'BY20BAPB30141000000000000002', amount: 1200 },
  { to_account: 'BY20BAPB30141000000000000003', amount: 2000 },
  { to_account: 'BY20BAPB30141000000000000004', amount: 3500 },
  { to_account: 'BY20BAPB30141000000000000005', amount: 700 },
  { to_account: 'BY20BAPB30141000000000000006', amount: 1500 },
  { to_account: 'BY20BAPB30141000000000000007', amount: 900 },
  { to_account: 'BY13PJCB30141000000000000000', amount: 500 },
  { to_account: 'BY13PJCB30141000000000000001', amount: 750 },
  { to_account: 'BY13PJCB30141000000000000002', amount: 1100 },
  { to_account: 'BY13PJCB30141000000000000003', amount: 1600 },
  { to_account: 'BY13PJCB30141000000000000004', amount: 2300 },
  { to_account: 'BY13PJCB30141000000000000005', amount: 3000 },
  { to_account: 'BY13PJCB30141000000000000006', amount: 4500 },
];

export default function () {
  const user = users[Math.floor(Math.random() * users.length)];

  const loginRes = http.post(
    `${BASE_URL}/login`,
    {
      email: user.email,
      password: user.password,
    }
  );

  check(loginRes, {
    'login ok (200/303)': (r) => r.status === 200 || r.status === 303,
  });


  const tx = transfers[Math.floor(Math.random() * transfers.length)];

  const txRes = http.post(
    `${BASE_URL}/transaction`,
    {
      amount: String(tx.amount),
      to_account: tx.to_account,
    }
  );

  check(txRes, {
    'transaction ok (200/303)': (r) => r.status === 200 || r.status === 303,
  });

  sleep(0.3);
}
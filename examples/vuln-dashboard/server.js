const express = require('express');
const app = express();
const PORT = process.env.PORT || 3000;

const vulnerabilities = [
  {
    id: 'CVE-2024-1234',
    severity: 'HIGH',
    package: 'openssl',
    affectedVersions: '< 3.0.8',
    status: 'fixed',
  },
  {
    id: 'CVE-2024-5678',
    severity: 'MEDIUM',
    package: 'curl',
    affectedVersions: '< 8.4.0',
    status: 'fixed',
  },
  {
    id: 'CVE-2024-9012',
    severity: 'LOW',
    package: 'zlib',
    affectedVersions: '< 1.3.1',
    status: 'investigating',
  },
];

app.get('/', (req, res) => {
  res.json({ service: 'vuln-dashboard', status: 'ok', version: '1.0.0' });
});

app.get('/vulns', (req, res) => {
  res.json(vulnerabilities);
});

app.get('/vulns/:id', (req, res) => {
  const vuln = vulnerabilities.find(v => v.id === req.params.id);
  if (!vuln) return res.status(404).json({ error: 'not found' });
  res.json(vuln);
});

app.listen(PORT, () => {
  console.log(`vuln-dashboard listening on :${PORT}`);
});

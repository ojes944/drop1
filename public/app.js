// Hauling App SPA logic
// Assumes backend at http://localhost:8080

const API_URL = 'http://localhost:8080/api';
let token = localStorage.getItem('token') || '';
let userName = '';
let map, userMarker;

function show(id) {
  document.querySelectorAll('#auth-forms form').forEach(f => f.classList.add('hidden'));
  document.getElementById(id).classList.remove('hidden');
}

function showMainApp() {
  document.getElementById('auth-forms').classList.add('hidden');
  document.getElementById('main-app').classList.remove('hidden');
  document.getElementById('user-name').textContent = userName;
  if (!map) initMap();
}

function showAuth() {
  document.getElementById('main-app').classList.add('hidden');
  document.getElementById('auth-forms').classList.remove('hidden');
}

// Auth Handlers
const signinForm = document.getElementById('signin-form');
signinForm.onsubmit = async e => {
  e.preventDefault();
  const email = document.getElementById('signin-email').value;
  const password = document.getElementById('signin-password').value;
  const res = await fetch(`${API_URL}/login`, {
    method: 'POST', headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({email, password})
  });
  if (res.ok) {
    const data = await res.json();
    token = data.token;
    localStorage.setItem('token', token);
    userName = email;
    showMainApp();
  } else {
    document.getElementById('signin-error').textContent = 'Invalid credentials.';
  }
};

document.getElementById('to-signup').onclick = () => show('signup-form');
document.getElementById('to-signin').onclick = () => show('signin-form');
document.getElementById('to-signin2').onclick = () => show('signin-form');
document.getElementById('to-reset').onclick = () => show('reset-form');

document.getElementById('signout').onclick = () => {
  token = '';
  localStorage.removeItem('token');
  showAuth();
};

// Signup
const signupForm = document.getElementById('signup-form');
signupForm.onsubmit = async e => {
  e.preventDefault();
  const name = document.getElementById('signup-name').value;
  const email = document.getElementById('signup-email').value;
  const password = document.getElementById('signup-password').value;
  const res = await fetch(`${API_URL}/signup`, {
    method: 'POST', headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({name, email, password})
  });
  if (res.ok) {
    userName = name;
    show('signin-form');
  } else {
    document.getElementById('signup-error').textContent = 'Signup failed.';
  }
};

// Reset Password
const resetForm = document.getElementById('reset-form');
resetForm.onsubmit = async e => {
  e.preventDefault();
  const email = document.getElementById('reset-email').value;
  const res = await fetch(`${API_URL}/reset_password`, {
    method: 'POST', headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({email})
  });
  if (res.ok) {
    document.getElementById('reset-success').textContent = 'Reset link sent (mock).';
    document.getElementById('reset-error').textContent = '';
  } else {
    document.getElementById('reset-error').textContent = 'Reset failed.';
    document.getElementById('reset-success').textContent = '';
  }
};

// MapLibre
function initMap() {
  map = new maplibregl.Map({
    container: 'map',
    style: 'https://demotiles.maplibre.org/style.json',
    center: [0, 0],
    zoom: 2
  });
}

// WebSocket logic (reuse ws_client.js if needed)
let ws;
function connectWS() {
  ws = new WebSocket('ws://localhost:8080/ws?token=' + encodeURIComponent(token));
  ws.onopen = () => {};
  ws.onmessage = e => {
    const msg = JSON.parse(e.data);
    if (msg.event === 'nearby_drivers') {
      showDrivers(msg.data);
    }
  };
  ws.onclose = () => setTimeout(connectWS, 2000);
}

function sendEvent(event, data) {
  if (ws && ws.readyState === 1) {
    ws.send(JSON.stringify({event, data}));
  }
}

document.getElementById('update-location').onclick = () => {
  if (!navigator.geolocation) return alert('Geolocation not supported');
  navigator.geolocation.getCurrentPosition(pos => {
    const {latitude, longitude} = pos.coords;
    sendEvent('update_location', {lat: latitude, lng: longitude});
    if (map) {
      map.setCenter([longitude, latitude]);
      if (userMarker) userMarker.remove();
      userMarker = new maplibregl.Marker({color: '#0078a8'}).setLngLat([longitude, latitude]).addTo(map);
    }
  });
};

document.getElementById('find-drivers').onclick = () => {
  if (!navigator.geolocation) return alert('Geolocation not supported');
  navigator.geolocation.getCurrentPosition(pos => {
    const {latitude, longitude} = pos.coords;
    sendEvent('get_nearby_drivers', {lat: latitude, lng: longitude, radius: 1000});
  });
};

function showDrivers(drivers) {
  const list = document.getElementById('drivers-list');
  list.innerHTML = '<h4>Nearby Drivers</h4>' + (drivers.length ? drivers.map(d => `<div>Driver ${d.id}: (${d.lat.toFixed(5)}, ${d.lng.toFixed(5)})</div>`).join('') : '<div>No drivers found.</div>');
  if (map) {
    drivers.forEach(d => {
      new maplibregl.Marker({color: '#c00'}).setLngLat([d.lng, d.lat]).addTo(map);
    });
  }
}

// Auto-login if token exists
if (token) {
  userName = localStorage.getItem('userName') || '';
  showMainApp();
  connectWS();
} else {
  showAuth();
}

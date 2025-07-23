// navigation helpers
function navigate(path) {
  history.pushState({}, '', path);
  render(path);
}

function render(path) {
  // скрыть все вьюшки
  document.querySelectorAll('#app > div').forEach(div => div.classList.add('hidden'));

  if (path === '/login') {
    document.getElementById('login-view').classList.remove('hidden');
  } else if (path === '/register') {
    document.getElementById('register-view').classList.remove('hidden');
  } else if (path === '/chat/' || path === '/chat') {
    document.getElementById('chat-view').classList.remove('hidden');
  } else {
    document.getElementById('main-view').classList.remove('hidden');
  }
}

async function checkAuthAndRedirect() {
  try {
    const res = await fetch('/api/me', {
      method: 'GET',
      credentials: 'include', // важно для кук!
    });
    if (!res.ok) return false;

    const user = await res.json();
    console.log('Auth user:', user);

    if (user && (user.un || user.username)) {
      history.pushState({}, '', '/chat/');
      render('/chat/');
      return true;
    }
  } catch (e) {
    console.error('Auth check failed:', e);
  }
  return false;
}

window.addEventListener('DOMContentLoaded', async () => {
  console.log('Page loaded');

  const redirected = await checkAuthAndRedirect();
  if (!redirected) {
    render(location.pathname);
  }

  // регистрация
  document.getElementById('register-form').onsubmit = async (e) => {
    e.preventDefault();
    const f = e.target;
    const res = await fetch('/api/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        pn: f.pn.value,
        un: f.un.value,
        email: f.email.value,
        password: f.password.value,
      })
    });

    if (res.ok) {
      navigate('/chat/');
    } else {
      const text = await res.text();
      alert('Registration failed: ' + text);
    }
  };

  // логин
  document.getElementById('login-form').onsubmit = async (e) => {
    e.preventDefault();
    const f = e.target;
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        email: f.email.value,
        password: f.password.value,
      })
    });

    if (res.ok) {
      navigate('/chat/');
    } else {
      const text = await res.text();
      alert('Login failed: ' + text);
    }
  };

  // переходы по истории
  window.onpopstate = () => {
    render(location.pathname);
  };
});

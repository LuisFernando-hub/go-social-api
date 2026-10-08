import http from 'k6/http';
import { check } from 'k6';

const BASE_URL = 'http://localhost:4001';

export const options = {
    vus: 1,
    iterations: 4,
};

const PASSWORD = 'Password123!';

export default function () {
    const id = __ITER;
    const email = `user_${id}@test.com`;

    const createUserResponse = http.post(
        `${BASE_URL}/api/v1/users`,
        JSON.stringify({
            username: `user_${id}`,
            email,
            password: PASSWORD,
        }),
        {
            headers: {
                'Content-Type': 'application/json',
            },
        }
    );

    const userAvailable = check(createUserResponse, {
        'usuário criado ou já existente': (r) =>
            r.status === 201 || r.status === 409,
    });

    if (!userAvailable) {
        console.error(
            `Falha no cadastro: ${createUserResponse.status} - ${createUserResponse.body}`
        );
        return;
    }

    if (createUserResponse.status === 201) {
        console.log(`Usuário criado: ${email}`);
    } else {
        console.log(`Usuário existente: ${email}`);
    }

    const loginResponse = http.post(
        `${BASE_URL}/api/v1/login`,
        JSON.stringify({
            email,
            password: PASSWORD,
        }),
        {
            headers: {
                'Content-Type': 'application/json',
            },
        }
    );

    const loginOk = check(loginResponse, {
        'login retornou 200': (r) => r.status === 200,
    });

    if (!loginOk) {
        console.error(
            `Falha no login: ${loginResponse.status} - ${loginResponse.body}`
        );
        return;
    }

    const loginData = loginResponse.json();
    const token = loginData.token;

    if (!token) {
        console.error('Token não encontrado na resposta do login');
        return;
    }

    const userId = loginData.user?.id

    if (userId === undefined || userId === null) {
        console.error(
            `ID do usuário não encontrado na resposta: ${loginResponse.body}`
        );
        return;
    }

    console.log(`VU ${__VU} | Usuário ${email} | ID ${userId}`);

    const authParams = {
        headers: {
            Authorization: `Bearer ${token}`,
        },
    };

    const postsResponse = http.get(
        `${BASE_URL}/api/v1/posts`,
        authParams
    );

    const postsOk = check(postsResponse, {
        'GET posts retornou 200': (r) => r.status === 200,
    });

    if (!postsOk) {
        console.error(
            `Falha ao buscar posts: ${postsResponse.status} - ${postsResponse.body}`
        );
        return;
    }

    const posts = postsResponse.json();

    if (!Array.isArray(posts)) {
        console.error(
            `A resposta de posts não é um array: ${postsResponse.body}`
        );
        return;
    }

    const createPostResponse = http.post(
        `${BASE_URL}/api/v1/posts`,
        JSON.stringify({
            content: `Post criado pelo ${email}`,
        }),
        {
            headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        }
    );

    check(createPostResponse, {
        'POST post retornou 201': (r) => r.status === 201,
    });

    const otherPosts = posts.filter(
        (post) => String(post.user_id) !== String(userId)
    );

    if (otherPosts.length === 0) {
        console.log(
            `VU ${__VU}: nenhum post de outro usuário disponível`
        );
        return;
    }

    const randomIndex = Math.floor(
        Math.random() * otherPosts.length
    );

    const otherPost = otherPosts[randomIndex];

    check(otherPosts, {
        'encontrou post de outro usuário': () =>
            otherPost !== undefined,
    });

    if (!otherPost) {
        return;
    }

    const type = Math.random() < 0.8
        ? 'LIKE'
        : 'DESLIKE';

    console.log(
        `Interação: ${type} | Usuário: ${userId} | Post: ${otherPost.id}`
    );

    const interactionResponse = http.post(
        `${BASE_URL}/api/v1/posts/interation`,
        JSON.stringify({
            post_id: otherPost.id,
            type,
        }),
        {
            headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        }
    );

    const interactionOk = check(interactionResponse, {
        'interação criada': (r) =>
            r.status === 200 || r.status === 201,
    });

    if (!interactionOk) {
        console.error(
            `Falha na interação: ${interactionResponse.status} - ${interactionResponse.body}`
        );
        return;
    }

    console.log(
        `Resposta da interação: ${interactionResponse.body}`
    );
}

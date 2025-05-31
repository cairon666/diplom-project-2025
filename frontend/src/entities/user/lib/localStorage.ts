export function setUserIdLocalStorage(userId: string) {
    localStorage.setItem('userId', userId);
}

export function getUserIdLocalStorage(): string | null {
    return localStorage.getItem('userId');
}

export function clearUserIdLocalStorage() {
    localStorage.removeItem('userId');
}

export function getAccessTokenLocalStorage(): string | null {
    return localStorage.getItem('accessToken');
}

export function setAccessTokenLocalStorage(accessToken: string) {
    localStorage.setItem('accessToken', accessToken);
}

export function clearAccessTokenLocalStorage() {
    localStorage.removeItem('accessToken');
}

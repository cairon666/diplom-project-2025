import { createRouter, RouterProvider } from '@tanstack/react-router';
import { Provider } from 'react-redux';
import { getUserIdLocalStorage } from 'src/entities/user/lib/localStorage';
import { store } from 'src/store';
import { Toaster } from 'ui/sonner';

import './styles/global.css';
import { routeTree } from './routeTree.gen';

export const router = createRouter({
    routeTree,
    context: {
        userId: getUserIdLocalStorage(),
    },
});

function App() {
    return (
        <Provider store={store}>
            <RouterProvider router={router} />
            <Toaster richColors />
        </Provider>
    );
}

export default App;

<script setup>
</script>

<template>
    <nav class="navbar">
        <ul class="navbar-nav">
            <li class="nav-item">
                <a href='#/' class="nav-link">
                    <span class="link-text">One</span>
                </a>
            </li>
            <li class="nav-item">
                <a href='#/' class="nav-link">
                    <span class="link-text">Two</span>
                </a>
            </li>
            <li class="nav-item">
                <a href='#/' class="nav-link">
                    <span class="link-text">Three</span>
                </a>
            </li>
        </ul>
    </nav>
</template>

<style scoped>
.navbar {
    position: fixed;
    z-index: 1000;
}

.navbar-nav {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
}

.nav-item {
    width: 100%
}

.nav-link {
    display: flex;
    align-items: center;
    height: 50x;
    color: red;
    text-decoration: none;
    filter: grayscale(100%) opacity(0.7)
}

.nav-link:hover {
    filter: grayscale(0%) opacity(1) brightness(100%);
    background-color: blue;
    color: white
}

.link-text {
    display: none;
    color: whitesmoke;
}

.nav-link svg {
    min-width: 2rem;
    margin: 0 1.5rem;
    fill: yellow;
    transition: 300ms;
}

@media only screen and (min-width: 600px) {
    .navbar {
        top: 0;
        width: 2rem;
        height: 100vh;
    }

    .navbar:hover {
        width: 16rem;
        background-color: #171717;
    }

    .navbar:hover .link-text {
        display: block;
    }

    main {
        margin-left: 5rem;
    }
}

</style>

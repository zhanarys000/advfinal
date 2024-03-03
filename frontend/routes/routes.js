const express = require('express');
const router = express.Router();
const axios = require('axios');
const { LocalStorage } = require('node-localstorage');
const localStorage = new LocalStorage('./scratch');


router.get('/', (req, res) => {
    console.log("start / ")
    res.status(200).render('register');
});

router.post('/register', async (req, res)=>{
    console.log("start register")
    const {email, username}= req.body
    axios.post('http://localhost:8080/api/getverificatoincode', { email })
    .then(response => {
        if (response.status == 200) {
            res.redirect(`/code?email=${email}&username=${username}`)
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
})
router.get('/code', async(req, res)=>{
    console.log("start code get ")
    const username = req.query.username
    const email = req.query.email
    res.render('code', {username, email})
})
router.post('/code', async(req, res)=>{
    console.log("start code post")

    const username = req.query.username
    const email = req.query.email
    const code = req.body.code
    axios.post('http://localhost:8080/api/checkverificatiocode', {
            email: email,
            code: code
    })
    .then(response => {
        if (response.status == 200) {
            res.redirect(`/getpass?email=${email}&username=${username}`)
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
})
router.get('/getpass', (req, res)=>{
    console.log("start getpass get")

    const username = req.query.username
    const email = req.query.email
    res.render('getpass', {username, email})
})
router.post('/getpass', (req, res)=>{
    console.log("start getpass post")

    const username = req.query.username
    const email = req.query.email
    const password = req.body.password
    axios.post('http://localhost:8080/api/register', {
            email: email,
            username: username, 
            password: password
        
    })
    .then(response => {
        if (response.status == 200) {
            const token = response.data.token;
            localStorage.setItem('token', token);
            res.redirect(`/login`)
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
})


router.get('/login', (req, res) => {
    console.log("start login get")

    res.status(200).render('login');
});

router.post('/login', async (req, res)=>{
    console.log("start login post")

    const {email, username, password}= req.body
    axios.post('http://localhost:8080/api/login', {  username, password })
    .then(response => {
        const token = response.data.token;
        localStorage.setItem('token', token);

        if (response.status == 200) {
            res.redirect(`index`)
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
})

router.get('/index', (req, res) => {
    console.log("start index");
    const token = localStorage.getItem('token');
    let { page, pageSize, sortBy, sortOrder, filterBy, filterValue } = req.query; 
    const currentPage = page || 1;
    const size = pageSize || 10;
    const sort = sortBy || "Name";
    const order = sortOrder || "asc";
    filterBy = filterBy || '';
    filterValue = filterValue || "";

    console.log(currentPage, size, sort, order, filterBy, filterValue)
    axios.get(`http://localhost:8080/api/books/all?page=${currentPage}&pageSize=${size}&sortBy=${sort}&sortOrder=${order}&filterBy=${filterBy}&filterValue=${filterValue}`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => {
        if (response.status == 200) {
            console.log(response.data)
            res.render('index', {books:response.data})
        }
    })
    .catch(error => {
        console.error(`Error fetching books for filter ${filterBy}:`, error);
    });
    
});


router.get("/buy/:bookId", (req, res) => {
    console.log("start enroll")

    const token = localStorage.getItem('token');
    const bookId = req.params.bookId;
    
    axios.post(`http://localhost:8080/api/books/buy/${bookId}`, {},{
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => {
        if (response.status === 200) {
            res.redirect(`/mybooks`);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
});

router.get('/profile', (req, res)=>{
    console.log("start profile")
    const token = localStorage.getItem('token');
    axios.get('http://localhost:8080/api/profile', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => {
        if (response.status == 200) {
            console.log(response.data.data)
            res.render(`profile`, {user: response.data.data})
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
})

router.get('/logout', (req, res)=>{
    res.redirect("/")
})

router.post('/update', (req, res) => {
    try {
        console.log("start update user get")
        const token = localStorage.getItem('token');
        const { username, email } = req.body;
        axios.put('http://localhost:8080/api/update',{username, email}, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        })
        .then(response => {
            if (response.status == 200) {
                res.redirect(`/profile`)
            }
        })
        .catch(error => {
            console.error('Error:', error);
            res.status(500).send('Internal Server Error');
        });
    } catch (error) {
        console.error('Error updating user information:', error);
        res.status(500).send('Internal Server Error');
    }
});

router.get('/subscribe', (req, res)=>{
    console.log("start subscribe get")
    const token = localStorage.getItem('token');
    axios.post('http://localhost:8080/api/subscribe', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => {
        if (response.status == 200 || response.status == 400) {
            console.log(response.data)
            res.redirect('/index')
        }
      
    })
    .catch(error => {
        console.error('Error:', error);
    });
})

router.post('/admin/spam', (req, res)=>{
    console.log("start subscribe get")
    const text = req.body.text
    console.log(text)
    const token = localStorage.getItem('token');
    axios.post('http://localhost:8080/api/send',{text},{
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => {
        if (response.status == 200 || response.status == 400) {
            console.log(response.data)
            res.redirect('/admin')
        }
      
    })
    .catch(error => {
        console.error('Error:', error);
    });
})
router.get('/mybooks', (req,res)=>{
    console.log("start mybooks get")
    const token = localStorage.getItem('token');
    axios.get('http://localhost:8080/api/books/my', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => {
        if (response.status == 200) {
            console.log(response.data)
            res.render(`mybooks`, {books: response.data.books})
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
})
router.get('/admin', (req, res)=>{
    console.log("start admin get")
    const token = localStorage.getItem('token');
    axios.get('http://localhost:8080/api/books', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => {
        if (response.status == 200) {
            console.log( response.data[1])
            res.render(`admin`, {books: response.data})
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
})


router.post("/admin", (req, res) => {
    console.log("start admin create")

    const token = localStorage.getItem('token');
    axios.post('http://localhost:8080/api/books',{
        name: req.body.name,
        description: req.body.description,
        price: req.body.price,
        genre: req.body.genre,
        isbn: req.body.isbn
    }, {
        headers: {
            'Authorization': `Bearer ${token}`, 
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        if (response.status == 201) {
            res.redirect(`/admin`)
        }
    })
    .catch(error => {
        console.error('Error:', error);
        res.status(500).json({ error: 'Ошибка сервера' });
    });
});

router.get('/admin/delete/:bookId', (req, res)=>{
    console.log("start admin dleete post")

    const token = localStorage.getItem('token');
    const bookId = req.params.bookId
    axios.delete(`http://localhost:8080/api/books/${bookId}`, {
        headers: {
            'Authorization': `Bearer ${token}`, 
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        if (response.status == 200) {
            res.redirect(`/admin`)
        }
    })
    .catch(error => {
        console.error('Error:', error);
        res.status(500).json({ error: 'Ошибка сервера' });
    });
})

router.post("/admin/update/:bookId", (req, res) => {
    console.log("start admin update");

    const token = localStorage.getItem('token');
    const bookId = req.params.bookId;
    axios.put(`http://localhost:8080/api/books/${bookId}`, {
        name: req.body.name,
        description: req.body.description,
        price: req.body.price,
        genre: req.body.genre,
        isbn: req.body.isbn
    }, {
        headers: {
            'Authorization': `Bearer ${token}`,
        }
    })
    .then(response => {
        if (response.status == 200) {
            res.redirect(`/admin`);
        }
    })
    .catch(error => {
        console.error('Error:', error);
        res.status(500).json({ error: 'Ошибка сервера' });
    });
});



module.exports = router;


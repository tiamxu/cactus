document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById('search');
    const categories = document.querySelectorAll('.category');
    
    searchInput.addEventListener('input', function() {
        const searchTerm = this.value.toLowerCase();
        
        categories.forEach(category => {
            const links = category.querySelectorAll('.link');
            let hasVisibleLinks = false;
            
            links.forEach(link => {
                const title = link.querySelector('span').textContent.toLowerCase();
                if (title.includes(searchTerm)) {
                    link.style.display = 'flex';
                    hasVisibleLinks = true;
                } else {
                    link.style.display = 'none';
                }
            });
            
            category.style.display = hasVisibleLinks ? 'block' : 'none';
        });
    });
});
document.addEventListener('DOMContentLoaded', () => {
    fetchProperties();
    fetchCustomers();
    
    // Initialize icons if lucide is available
    if (window.lucide) {
        lucide.createIcons();
    }

    const addForm = document.getElementById('add-property-form');
    if (addForm) {
        addForm.addEventListener('submit', handleAddProperty);
    }
    
    const addCxForm = document.getElementById('add-customer-form');
    if (addCxForm) {
        addCxForm.addEventListener('submit', handleAddCustomer);
    }
});

// View Switching Logic
window.switchView = function(viewName, element) {
    // Hide all views
    document.querySelectorAll('.view-section').forEach(el => {
        el.style.display = 'none';
    });
    
    // Show selected view
    const view = document.getElementById(`view-${viewName}`);
    if (view) {
        view.style.display = 'block';
    }
    
    // Update active state in nav
    document.querySelectorAll('.nav-links a').forEach(el => {
        el.classList.remove('active');
    });
    if (element) {
        element.classList.add('active');
    }

    // Refresh data if needed
    if (window.lucide) lucide.createIcons();
}

// --- PROPERTIES LOGIC ---

async function fetchProperties() {
    try {
        const response = await fetch('/api/v1/properties');
        if (!response.ok) throw new Error('Failed to fetch properties');
        
        const data = await response.json();
        renderProperties(data);
        updateStats(data, 'properties');
    } catch (error) {
        console.error('Error fetching properties:', error);
    }
}

function renderProperties(properties) {
    const tableBody = document.getElementById('property-table-body');
    if (!tableBody) return;
    
    tableBody.innerHTML = '';

    if (!properties || properties.length === 0) {
        tableBody.innerHTML = `<tr><td colspan="5" style="text-align: center; color: var(--text-secondary); padding: 2rem;">No properties found.</td></tr>`;
        return;
    }

    properties.forEach(property => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>
                <div style="font-weight: 600;">${property.title}</div>
                <div style="font-size: 0.75rem; color: var(--text-secondary);">${new Date(property.created_at).toLocaleDateString()}</div>
            </td>
            <td>${property.address}</td>
            <td style="font-family: monospace; font-size: 0.9rem;">${formatCurrency(property.price)}</td>
            <td><span class="status-badge active">Active</span></td>
            <td>
                <button class="action-btn" onclick="deleteProperty(${property.id})" title="Delete"><i data-lucide="trash-2"></i></button>
            </td>
        `;
        tableBody.appendChild(row);
    });
    
    if (window.lucide) lucide.createIcons();
}

async function handleAddProperty(e) {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const data = {
        title: formData.get('title'),
        address: formData.get('address'),
        description: formData.get('description'),
        price: parseFloat(formData.get('price'))
    };

    try {
        const response = await fetch('/api/v1/properties', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (!response.ok) throw new Error('Failed to create property');

        await fetchProperties();
        closeAddModal();
        e.target.reset();
    } catch (error) {
        console.error(error);
        alert('Failed: ' + error.message);
    }
}

async function deleteProperty(id) {
    if (!confirm('Delete this property?')) return;
    try {
        await fetch(`/api/v1/properties/${id}`, { method: 'DELETE' });
        await fetchProperties();
    } catch (error) {
        console.error(error);
        alert('Failed to delete');
    }
}

// --- CUSTOMERS LOGIC ---

async function fetchCustomers() {
    try {
        const response = await fetch('/api/v1/customers');
        if (!response.ok) throw new Error('Failed to fetch customers');
        
        const data = await response.json();
        renderCustomers(data);
        updateStats(data, 'customers');
    } catch (error) {
        console.error('Error fetching customers:', error);
    }
}

function renderCustomers(customers) {
    const tableBody = document.getElementById('customer-table-body');
    if (!tableBody) return;
    
    tableBody.innerHTML = '';

    if (!customers || customers.length === 0) {
        tableBody.innerHTML = `<tr><td colspan="5" style="text-align: center; color: var(--text-secondary); padding: 2rem;">No customers found.</td></tr>`;
        return;
    }

    customers.forEach(c => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>
                <div style="font-weight: 600;">${c.name}</div>
                <div style="font-size: 0.75rem; color: var(--text-secondary);">ID: ${c.id}</div>
            </td>
            <td>${c.email}</td>
            <td>${c.phone || '-'}</td>
            <td><span class="status-badge ${c.status === 'Active' ? 'active' : ''}">${c.status}</span></td>
            <td>
                <button class="action-btn" onclick="deleteCustomer(${c.id})" title="Delete"><i data-lucide="trash-2"></i></button>
            </td>
        `;
        tableBody.appendChild(row);
    });
    
    if (window.lucide) lucide.createIcons();
}

async function handleAddCustomer(e) {
    e.preventDefault();
    const formData = new FormData(e.target);
    const data = {
        name: formData.get('name'),
        email: formData.get('email'),
        phone: formData.get('phone'),
        status: 'Active'
    };

    try {
        const response = await fetch('/api/v1/customers', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (!response.ok) throw new Error('Failed to create customer');

        await fetchCustomers();
        closeAddCustomerModal();
        e.target.reset();
    } catch (error) {
        console.error(error);
        alert('Failed: ' + error.message);
    }
}

async function deleteCustomer(id) {
    if (!confirm('Delete this customer?')) return;
    try {
        await fetch(`/api/v1/customers/${id}`, { method: 'DELETE' });
        await fetchCustomers();
    } catch (error) {
        console.error(error);
        alert('Failed to delete');
    }
}

// --- UTILS ---

function updateStats(data, type) {
    if (type === 'properties') {
        const el = document.getElementById('total-properties');
        if (el) el.textContent = data ? data.length : 0;
    } else if (type === 'customers') {
        const el = document.getElementById('total-customers');
        if (el) el.textContent = data ? data.length : 0;
    }
}

function formatCurrency(amount) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(amount);
}

// --- MODALS ---

window.showAddModal = function() {
    const modal = document.getElementById('add-modal');
    if (modal) modal.style.display = 'block';
}

window.closeAddModal = function() {
    const modal = document.getElementById('add-modal');
    if (modal) modal.style.display = 'none';
}

window.showAddCustomerModal = function() {
    const modal = document.getElementById('add-customer-modal');
    if (modal) modal.style.display = 'block';
}

window.closeAddCustomerModal = function() {
    const modal = document.getElementById('add-customer-modal');
    if (modal) modal.style.display = 'none';
}

// Close modal when clicking outside
window.onclick = function(event) {
    const propModal = document.getElementById('add-modal');
    const custModal = document.getElementById('add-customer-modal');
    if (event.target == propModal) {
        propModal.style.display = "none";
    }
    if (event.target == custModal) {
        custModal.style.display = "none";
    }
}

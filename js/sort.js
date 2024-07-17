function sortByDate(elements, order) {
  const elementsArray = Array.from(elements);
  const parentContainer = elementsArray[0].parentNode;

  elementsArray.sort((a, b) => {
    const dateA = new Date(a.dataset.date);
    const dateB = new Date(b.dataset.date);

    if (order === 'asc') {
      return dateA - dateB;
    } else {
      return dateB - dateA;
    }
  });

  elementsArray.forEach((element) => {
    parentContainer.appendChild(element);
  });

  // Toggle order for the next click
  const nextOrder = order === 'asc' ? 'desc' : 'asc';

  // Update sorting icon based on the current order
  updateSortingIcon(order);

  // Attach click event listener to toggle order on the next click
  parentContainer.addEventListener('click', () => {
    sortByDate(elements, nextOrder);
  });
}

function updateSortingIcon(order) {
  // You can add logic here to update sorting icon based on the order
  // For example, toggle a CSS class or update an <i> element
  const sortingIcon = document.getElementById('sorting-icon');
  if (sortingIcon) {
    sortingIcon.innerText = order === 'asc' ? '🔼' : '🔽';
  }
}
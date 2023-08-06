const studentElements = document.querySelectorAll('.directory-Entry');

// Array to store the extracted student data
const students = [];

// Loop through each student element
studentElements.forEach((studentElement) => {
    // Object to store the student data
    const student = {};

    // Extract the student's name
    const nameElement = studentElement.querySelector('.directory-Entry_Title');
    if (nameElement) {
        student.name = nameElement.textContent.trim();
    }

    // Extract the student's email
    const emailElement = studentElement.querySelector('.directory-Entry_FieldValue a[href^="mailto:"]');
    if (emailElement) {
        student.email = emailElement.textContent.trim();
    }

    // Extract the student's parent's email
    const parentEmailElement = studentElement.querySelector('.directory-Entry_HouseholdSection .directory-Entry_FieldValue a[href^="mailto:"]');
    if (parentEmailElement) {
        student.parentEmail = parentEmailElement.textContent.trim();
    }

    // Extract the state where the student lives
    const addressElement = studentElement.querySelector('.directory-Entry_HouseholdSection .directory-Entry_FieldTitle');
    if (addressElement) {
        const addressParts = addressElement.textContent.trim().split(',');
        if (addressParts.length > 1) {
            student.state = addressParts[addressParts.length - 2].trim();
        }
    }

    // Add the student data to the array
    students.push(student);
});
console.log(students)
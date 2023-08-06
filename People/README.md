# People Directory

Welcome to the People directory for Avon enrollment. This directory contains information about all the people who are currently enrolled in Avon. Additionally, it includes JavaScript files that can be used to update the information of individuals.

## Table of Contents

- [Directory Structure](#directory-structure)
- [Usage](#usage)
- [Updating People](#updating-people)

## Directory Structure

The directory structure is organized as follows:\
├── People/\
│ ├── js/\
│ │ ├── updateStudents.js\
│ │ └── updateTeachers.js\
│ ├── students.json\
│ ├── teachers.json\
└── ...

- The `js/` directory contains JavaScript files used to update the information of people in Avon.
- Each student enrolled in Avon is represented in a JSON file (`students.json`) within the `People/` directory.
- Each teacher enrolled in Avon is represented in a JSON file (`teachers.json`) within the `People/` directory.

## Usage

To access the information about a particular person, locate their corresponding JSON file in the `People/` directory. Each JSON file contains details about the individual, such as their name, contact information, and enrollment status.

Feel free to explore the JSON files and retrieve the necessary information based on your requirements.

## Updating People

To update the information of individuals in the Avon enrollment, you can utilize the provided JavaScript files in the `js/` directory. These files are specifically designed to handle the updating process.

For example, you can use the `updatePerson.js` file to modify the details of a person within the system. Refer to the documentation within the JavaScript files for instructions on how to use them effectively.

Please exercise caution when updating the information and ensure that any changes made are accurate and necessary.

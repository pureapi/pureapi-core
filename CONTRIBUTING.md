# Contributing to PureAPI Core

Thank you for your interest in contributing to PureAPI Core! Your contributions help us maintain a robust, lightweight framework for building web APIs in Go.


## How to Contribute

### Reporting Issues
- **Bug Reports:**  
  Please report bugs using our GitHub issue tracker. Provide a clear description, reproduction steps, and any relevant logs or error messages.
- **Feature Requests:**  
  If you have an enhancement idea, open an issue or discussion with a detailed explanation of the proposed improvement.

### Submitting Pull Requests
1. **Fork and Clone:**  
   Fork the repository and clone your fork locally.
2. **Create a Branch:**  
   Use a descriptive branch name (e.g., `feature/add-authentication` or `bugfix/fix-db-connection`).
3. **Implement Your Changes:**  
   Adhere to our coding guidelines and ensure your modifications are covered by tests according to our testing strategy:
   - **Public Contracts:** Ensure all externally visible APIs work as intended.
   - **Security-Critical Areas:** Include additional tests for sensitive code paths.
   - **Private Functions:** Test indirectly via the public interface unless direct testing is essential.
4. **Commit and Push:**  
   Write clear, concise commit messages and push your changes to your fork.
5. **Create a Pull Request:**  
   Open a pull request describing your changes, linking to any relevant issues, and outlining your testing approach.

### Coding Guidelines
- Follow Go best practices and maintain consistency in coding style.
- Write clear, well-documented code with meaningful names.
- Ensure that new code is accompanied by appropriate tests.

---

## Documentation

### Docs Directory Structure

We strive to maintain a minimal yet clear set of documentation to help everyone understand the system.

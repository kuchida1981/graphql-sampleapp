# Proposal: Refactor Repository Structure

This proposal addresses the unstructured layout of the `internal/repository` directory. Currently, repository implementations for different data stores (Postgres, Firestore) and their corresponding interfaces are mixed together, leading to a confusing and disorganized structure.

This change introduces a more organized, layered structure by moving data store-specific implementations into their own subdirectories. This will improve code navigation, clarify dependencies, and establish a clear pattern for adding new data stores in the future.

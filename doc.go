// Package ece provides the low-level data structures of D&Gine's
// Entity Component Event (ECE) architecture.
//
// ECE is a hand-made variant of the classic Entity Component System (ECS)
// where the "System" part is removed: instead of systems iterating over
// components every frame, behaviour is driven exclusively by Events.
// Entities are plain integer IDs, and components are stored in
// SparseArray containers (one per component type).
package ece

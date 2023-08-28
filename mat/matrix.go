package mat

// Vector represents a vector with an associated element increment.
type Vector[T any] struct {
	N    int
	Data []float64
	Inc  int
}

// Matrix is the basic matrix interface type.
type Matrix[T any] interface {
	// Dims returns the dimensions of a Matrix.
	Dims() (r, c int)

	// At returns the value of a matrix element at row i, column j.
	// It will panic if i or j are out of bounds for the matrix.
	At(i, j int) T

	// T returns the transpose of the Matrix. Whether T returns a copy of the
	// underlying data is implementation dependent.
	// This method may be implemented using the Transpose type, which
	// provides an implicit matrix transpose.
	T() Matrix[T]
}

// Transpose is a type for performing an implicit matrix transpose. It implements
// the Matrix interface, returning values from the transpose of the matrix within.
type Transpose[T any] struct {
	Matrix Matrix[T]
}

// At returns the value of the element at row i and column j of the transposed
// matrix, that is, row j and column i of the Matrix field.
func (t Transpose[T]) At(i, j int) T {
	return t.Matrix.At(j, i)
}

// Dims returns the dimensions of the transposed matrix. The number of rows returned
// is the number of columns in the Matrix field, and the number of columns is
// the number of rows in the Matrix field.
func (t Transpose[T]) Dims() (r, c int) {
	c, r = t.Matrix.Dims()
	return r, c
}

// T performs an implicit transpose by returning the Matrix field.
func (t Transpose[T]) T() Matrix[T] {
	return t.Matrix
}

// Untranspose returns the Matrix field.
func (t Transpose[T]) Untranspose() Matrix[T] {
	return t.Matrix
}

type matrix[T any] struct {
	Rows, Cols int
	Data       []T
	Stride     int

	capRows, capCols int
}

// New creates a new matrix with r rows and c columns. If data == nil,
// a new slice is allocated for the backing slice. If len(data) == r*c, data is
// used as the backing slice, and changes to the elements of the returned Dense
// will be reflected in data. If neither of these is true, NewMatrix will panic.
// NewDense will panic if either r or c is zero.
//
// The data must be arranged in row-major order, i.e. the (i*c + j)-th
// element in the data slice is the {i, j}-th element in the matrix.
func New[T any](r, c int, data []T) Matrix[T] {
	if r <= 0 || c <= 0 {
		if r == 0 || c == 0 {
			panic(ErrZeroLength)
		}
		panic(ErrNegativeDimension)
	}
	if data != nil && r*c != len(data) {
		panic(ErrShape)
	}
	if data == nil {
		data = make([]T, r*c)
	}
	return &matrix[T]{
		Rows:    r,
		Cols:    c,
		Stride:  c,
		Data:    data,
		capRows: r,
		capCols: c,
	}
}

// Dims returns the number of rows and columns in the matrix.
func (m *matrix[T]) Dims() (r, c int) { return m.Rows, m.Cols }

// Caps returns the number of rows and columns in the backing matrix.
func (m *matrix[T]) Caps() (r, c int) { return m.capRows, m.capCols }

func (m *matrix[T]) At(i, j int) T {
	return *new(T)
}

// T performs an implicit transpose by returning the receiver inside a Transpose.
func (m *matrix[T]) T() Matrix[T] {
	return Transpose[T]{m}
}

// Error represents matrix handling errors. These errors can be recovered by Maybe wrappers.
type Error struct{ string }

func (err Error) Error() string { return err.string }

var (
	ErrNegativeDimension   = Error{"mat: negative dimension"}
	ErrIndexOutOfRange     = Error{"mat: index out of range"}
	ErrReuseNonEmpty       = Error{"mat: reuse of non-empty matrix"}
	ErrRowAccess           = Error{"mat: row index out of range"}
	ErrColAccess           = Error{"mat: column index out of range"}
	ErrVectorAccess        = Error{"mat: vector index out of range"}
	ErrZeroLength          = Error{"mat: zero length in matrix dimension"}
	ErrRowLength           = Error{"mat: row length mismatch"}
	ErrColLength           = Error{"mat: col length mismatch"}
	ErrSquare              = Error{"mat: expect square matrix"}
	ErrNormOrder           = Error{"mat: invalid norm order for matrix"}
	ErrSingular            = Error{"mat: matrix is singular"}
	ErrShape               = Error{"mat: dimension mismatch"}
	ErrIllegalStride       = Error{"mat: illegal stride"}
	ErrPivot               = Error{"mat: malformed pivot list"}
	ErrTriangle            = Error{"mat: triangular storage mismatch"}
	ErrTriangleSet         = Error{"mat: triangular set out of bounds"}
	ErrBandwidth           = Error{"mat: bandwidth out of range"}
	ErrBandSet             = Error{"mat: band set out of bounds"}
	ErrDiagSet             = Error{"mat: diagonal set out of bounds"}
	ErrSliceLengthMismatch = Error{"mat: input slice length mismatch"}
	ErrNotPSD              = Error{"mat: input not positive symmetric definite"}
	ErrFailedEigen         = Error{"mat: eigendecomposition not successful"}
)

// ErrorStack represents matrix handling errors that have been recovered by Maybe wrappers.
type ErrorStack struct {
	Err error

	// StackTrace is the stack trace
	// recovered by Maybe, MaybeFloat
	// or MaybeComplex.
	StackTrace string
}

func (err ErrorStack) Error() string { return err.Err.Error() }

const badCap = "mat: bad capacity"

package mat

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

type General struct {
	Rows, Cols int
	Data       []float64
	Stride     int
}

type Matrix struct {
	mat              General
	capRows, capCols int
}

// New creates a new matrix with r rows and c columns. If data == nil,
// a new slice is allocated for the backing slice. If len(data) == r*c, data is
// used as the backing slice, and changes to the elements of the returned Dense
// will be reflected in data. If neither of these is true, New will panic.
// New will panic if either r or c is zero.
//
// The data must be arranged in row-major order, i.e. the (i*c + j)-th
// element in the data slice is the {i, j}-th element in the matrix.
func New(r, c int, data []float64) *Matrix {
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
		data = make([]float64, r*c)
	}
	return &Matrix{
		mat: General{
			Rows:   r,
			Cols:   c,
			Stride: c,
			Data:   data,
		},
		capRows: r,
		capCols: c,
	}
}

/*
type Matrix[T any] struct {
	buff [][]T
}
*/

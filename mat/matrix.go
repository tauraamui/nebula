package mat

import "unsafe"

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
// used as the backing slice, and changes to the elements of the returned Matrix
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

// ReuseAs changes the receiver if it IsEmpty() to be of size r×c.
//
// ReuseAs re-uses the backing data slice if it has sufficient capacity,
// otherwise a new slice is allocated. The backing data is zero on return.
//
// ReuseAs panics if the receiver is not empty, and panics if
// the input sizes are less than one. To empty the receiver for re-use,
// Reset should be used.
func (m *Matrix) ReuseAs(r, c int) {
	if r <= 0 || c <= 0 {
		if r == 0 || c == 0 {
			panic(ErrZeroLength)
		}
		panic(ErrNegativeDimension)
	}
	if !m.IsEmpty() {
		panic(ErrReuseNonEmpty)
	}
	m.reuseAsZeroed(r, c)
}

// reuseAsNonZeroed resizes an empty matrix to a r×c matrix,
// or checks that a non-empty matrix is r×c. It does not zero
// the data in the receiver.
func (m *Matrix) reuseAsNonZeroed(r, c int) {
	// reuseAs must be kept in sync with reuseAsZeroed.
	if m.mat.Rows > m.capRows || m.mat.Cols > m.capCols {
		// Panic as a string, not a mat.Error.
		panic(badCap)
	}
	if r == 0 || c == 0 {
		panic(ErrZeroLength)
	}
	if m.IsEmpty() {
		m.mat = blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: c,
			Data:   use(m.mat.Data, r*c),
		}
		m.capRows = r
		m.capCols = c
		return
	}
	if r != m.mat.Rows || c != m.mat.Cols {
		panic(ErrShape)
	}
}

// reuseAsZeroed resizes an empty matrix to a r×c matrix,
// or checks that a non-empty matrix is r×c. It zeroes
// all the elements of the matrix.
func (m *Matrix) reuseAsZeroed(r, c int) {
	// reuseAsZeroed must be kept in sync with reuseAsNonZeroed.
	if m.mat.Rows > m.capRows || m.mat.Cols > m.capCols {
		// Panic as a string, not a mat.Error.
		panic(badCap)
	}
	if r == 0 || c == 0 {
		panic(ErrZeroLength)
	}
	if m.IsEmpty() {
		m.mat = blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: c,
			Data:   useZeroed(m.mat.Data, r*c),
		}
		m.capRows = r
		m.capCols = c
		return
	}
	if r != m.mat.Rows || c != m.mat.Cols {
		panic(ErrShape)
	}
	m.Zero()
}

// Zero sets all of the matrix elements to zero.
func (m *Matrix) Zero() {
	r := m.mat.Rows
	c := m.mat.Cols
	for i := 0; i < r; i++ {
		zero(m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+c])
	}
}

// isolatedWorkspace returns a new dense matrix w with the size of a and
// returns a callback to defer which performs cleanup at the return of the call.
// This should be used when a method receiver is the same pointer as an input argument.
func (m *Matrix) isolatedWorkspace(a Matrix) (w *Matrix, restore func()) {
	r, c := a.Dims()
	if r == 0 || c == 0 {
		panic(ErrZeroLength)
	}
	w = getMatrixWorkspace(r, c, false)
	return w, func() {
		m.Copy(w)
		putMatrixWorkspace(w)
	}
}

// Reset empties the matrix so that it can be reused as the
// receiver of a dimensionally restricted operation.
//
// Reset should not be used when the matrix shares backing data.
// See the Reseter interface for more information.
func (m *Matrix) Reset() {
	// Row, Cols and Stride must be zeroed in unison.
	m.mat.Rows, m.mat.Cols, m.mat.Stride = 0, 0, 0
	m.capRows, m.capCols = 0, 0
	m.mat.Data = m.mat.Data[:0]
}

// IsEmpty returns whether the receiver is empty. Empty matrices can be the
// receiver for size-restricted operations. The receiver can be emptied using
// Reset.
func (m *Matrix) IsEmpty() bool {
	// It must be the case that m.Dims() returns
	// zeros in this case. See comment in Reset().
	return m.mat.Stride == 0
}

// asTriMatrix returns a TriMatrix with the given size and side. The backing data
// of the TriMatrix is the same as the receiver.
func (m *Matrix) asTriMatrix(n int, diag blas.Diag, uplo blas.Uplo) *TriMatrix {
	return &TriMatrix{
		mat: blas64.Triangular{
			N:      n,
			Stride: m.mat.Stride,
			Data:   m.mat.Data,
			Uplo:   uplo,
			Diag:   diag,
		},
		cap: n,
	}
}

// MatrixCopyOf returns a newly allocated copy of the elements of a.
func MatrixCopyOf(a Matrix) *Matrix {
	d := &Matrix{}
	d.CloneFrom(a)
	return d
}

// SetRawMatrix sets the underlying blas64.General used by the receiver.
// Changes to elements in the receiver following the call will be reflected
// in b.
func (m *Matrix) SetRawMatrix(b blas64.General) {
	m.capRows, m.capCols = b.Rows, b.Cols
	m.mat = b
}

// RawMatrix returns the underlying blas64.General used by the receiver.
// Changes to elements in the receiver following the call will be reflected
// in returned blas64.General.
func (m *Matrix) RawMatrix() blas64.General { return m.mat }

// Dims returns the number of rows and columns in the matrix.
func (m *Matrix) Dims() (r, c int) { return m.mat.Rows, m.mat.Cols }

// Caps returns the number of rows and columns in the backing matrix.
func (m *Matrix) Caps() (r, c int) { return m.capRows, m.capCols }

// T performs an implicit transpose by returning the receiver inside a Transpose.
func (m *Matrix) T() Matrix {
	return Transpose{m}
}

// ColView returns a Vector reflecting the column j, backed by the matrix data.
//
// See ColViewer for more information.
func (m *Matrix) ColView(j int) Vector {
	var v VecMatrix
	v.ColViewOf(m, j)
	return &v
}

// SetCol sets the values in the specified column of the matrix to the values
// in src. len(src) must equal the number of rows in the receiver.
func (m *Matrix) SetCol(j int, src []float64) {
	if j >= m.mat.Cols || j < 0 {
		panic(ErrColAccess)
	}
	if len(src) != m.mat.Rows {
		panic(ErrColLength)
	}

	blas64.Copy(
		blas64.Vector{N: m.mat.Rows, Inc: 1, Data: src},
		blas64.Vector{N: m.mat.Rows, Inc: m.mat.Stride, Data: m.mat.Data[j:]},
	)
}

// SetRow sets the values in the specified rows of the matrix to the values
// in src. len(src) must equal the number of columns in the receiver.
func (m *Matrix) SetRow(i int, src []float64) {
	if i >= m.mat.Rows || i < 0 {
		panic(ErrRowAccess)
	}
	if len(src) != m.mat.Cols {
		panic(ErrRowLength)
	}

	copy(m.rawRowView(i), src)
}

// RowView returns row i of the matrix data represented as a column vector,
// backed by the matrix data.
//
// See RowViewer for more information.
func (m *Matrix) RowView(i int) Vector {
	var v VecMatrix
	v.RowViewOf(m, i)
	return &v
}

// RawRowView returns a slice backed by the same array as backing the
// receiver.
func (m *Matrix) RawRowView(i int) []float64 {
	if i >= m.mat.Rows || i < 0 {
		panic(ErrRowAccess)
	}
	return m.rawRowView(i)
}

func (m *Matrix) rawRowView(i int) []float64 {
	return m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+m.mat.Cols]
}

// DiagView returns the diagonal as a matrix backed by the original data.
func (m *Matrix) DiagView() Diagonal {
	n := min(m.mat.Rows, m.mat.Cols)
	return &DiagMatrix{
		mat: blas64.Vector{
			N:    n,
			Inc:  m.mat.Stride + 1,
			Data: m.mat.Data[:(n-1)*m.mat.Stride+n],
		},
	}
}

// Slice returns a new Matrix that shares backing data with the receiver.
// The returned matrix starts at {i,j} of the receiver and extends k-i rows
// and l-j columns. The final row in the resulting matrix is k-1 and the
// final column is l-1.
// Slice panics with ErrIndexOutOfRange if the slice is outside the capacity
// of the receiver.
func (m *Matrix) Slice(i, k, j, l int) Matrix {
	return m.slice(i, k, j, l)
}

func (m *Matrix) slice(i, k, j, l int) *Matrix {
	mr, mc := m.Caps()
	if i < 0 || mr <= i || j < 0 || mc <= j || k < i || mr < k || l < j || mc < l {
		if i == k || j == l {
			panic(ErrZeroLength)
		}
		panic(ErrIndexOutOfRange)
	}
	t := *m
	t.mat.Data = t.mat.Data[i*t.mat.Stride+j : (k-1)*t.mat.Stride+l]
	t.mat.Rows = k - i
	t.mat.Cols = l - j
	t.capRows -= i
	t.capCols -= j
	return &t
}

// Grow returns the receiver expanded by r rows and c columns. If the dimensions
// of the expanded matrix are outside the capacities of the receiver a new
// allocation is made, otherwise not. Note the receiver itself is not modified
// during the call to Grow.
func (m *Matrix) Grow(r, c int) Matrix {
	if r < 0 || c < 0 {
		panic(ErrIndexOutOfRange)
	}
	if r == 0 && c == 0 {
		return m
	}

	r += m.mat.Rows
	c += m.mat.Cols

	var t Matrix
	switch {
	case m.mat.Rows == 0 || m.mat.Cols == 0:
		t.mat = blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: c,
			// We zero because we don't know how the matrix will be used.
			// In other places, the mat is immediately filled with a result;
			// this is not the case here.
			Data: useZeroed(m.mat.Data, r*c),
		}
	case r > m.capRows || c > m.capCols:
		cr := max(r, m.capRows)
		cc := max(c, m.capCols)
		t.mat = blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: cc,
			Data:   make([]float64, cr*cc),
		}
		t.capRows = cr
		t.capCols = cc
		// Copy the complete matrix over to the new matrix.
		// Including elements not currently visible. Use a temporary structure
		// to avoid modifying the receiver.
		var tmp Matrix
		tmp.mat = blas64.General{
			Rows:   m.mat.Rows,
			Cols:   m.mat.Cols,
			Stride: m.mat.Stride,
			Data:   m.mat.Data,
		}
		tmp.capRows = m.capRows
		tmp.capCols = m.capCols
		t.Copy(&tmp)
		return &t
	default:
		t.mat = blas64.General{
			Data:   m.mat.Data[:(r-1)*m.mat.Stride+c],
			Rows:   r,
			Cols:   c,
			Stride: m.mat.Stride,
		}
	}
	t.capRows = r
	t.capCols = c
	return &t
}

// CloneFrom makes a copy of a into the receiver, overwriting the previous value of
// the receiver. The clone from operation does not make any restriction on shape and
// will not cause shadowing.
//
// See the ClonerFrom interface for more information.
func (m *Matrix) CloneFrom(a Matrix) {
	r, c := a.Dims()
	mat := General{
		Rows:   r,
		Cols:   c,
		Stride: c,
	}
	m.capRows, m.capCols = r, c

	amat := a.mat
	mat.Data = make([]float64, r*c)
	for i := 0; i < r; i++ {
		copy(mat.Data[i*c:(i+1)*c], amat.Data[i*amat.Stride:i*amat.Stride+c])
	}
	m.mat = mat
}

// Copy makes a copy of elements of a into the receiver. It is similar to the
// built-in copy; it copies as much as the overlap between the two matrices and
// returns the number of rows and columns it copied. If a aliases the receiver
// and is a transposed Matrix or VecMatrix, with a non-unitary increment, Copy will
// panic.
//
// See the Copier interface for more information.
func (m *Matrix) Copy(a *Matrix) (r, c int) {
	r, c = a.Dims()
	if a == m {
		return r, c
	}
	r = min(r, m.mat.Rows)
	c = min(c, m.mat.Cols)
	if r == 0 || c == 0 {
		return 0, 0
	}

	amat := a.mat
	switch o := offset(m.mat.Data, amat.Data); {
	case o < 0:
		for i := r - 1; i >= 0; i-- {
			copy(m.mat.Data[i*m.mat.Stride:i*m.mat.Stride+c], amat.Data[i*amat.Stride:i*amat.Stride+c])
		}
	case o > 0:
		for i := 0; i < r; i++ {
			copy(m.mat.Data[i*m.mat.Stride:i*m.mat.Stride+c], amat.Data[i*amat.Stride:i*amat.Stride+c])
		}
	default:
		// Nothing to do.
	}

	return r, c
}

// Stack appends the rows of b onto the rows of a, placing the result into the
// receiver with b placed in the greater indexed rows. Stack will panic if the
// two input matrices do not have the same number of columns or the constructed
// stacked matrix is not the same shape as the receiver.
func (m *Matrix) Stack(a, b *Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ac != bc || m == a || m == b {
		panic(ErrShape)
	}

	m.reuseAsNonZeroed(ar+br, ac)

	m.Copy(a)
	w := m.slice(ar, ar+br, 0, bc)
	w.Copy(b)
}

// Augment creates the augmented matrix of a and b, where b is placed in the
// greater indexed columns. Augment will panic if the two input matrices do
// not have the same number of rows or the constructed augmented matrix is
// not the same shape as the receiver.
func (m *Matrix) Augment(a, b *Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || m == a || m == b {
		panic(ErrShape)
	}

	m.reuseAsNonZeroed(ar, ac+bc)

	m.Copy(a)
	w := m.slice(0, br, ac, ac+bc)
	w.Copy(b)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// offset returns the number of float64 values b[0] is after a[0].
func offset(a, b []float64) int {
	if &a[0] == &b[0] {
		return 0
	}
	// This expression must be atomic with respect to GC moves.
	// At this stage this is true, because the GC does not
	// move. See https://golang.org/issue/12445.
	return int(uintptr(unsafe.Pointer(&b[0]))-uintptr(unsafe.Pointer(&a[0]))) / int(unsafe.Sizeof(float64(0)))
}

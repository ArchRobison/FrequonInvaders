# Julia program for generating Go assembly language x86 kernels for Fourier Invaders
# The program is *not* a general purpose assembly code generator.
# It contains just enough to generate the kernals of interest.

# Floating-point operand (scalar or vector)
abstract Operand

# Memory location
abstract Mem <: Operand

const OperandOrImmediate = Union{Operand,Integer,Float32}

asm(i::Int) = "\$$(string(i))"
asm(i::Float32) = "\$$(string(i))"
asm(u::UInt8) = u<10 ? "\$$u" : "\$0x$(hex(u))"

#---------------------------------------
# Vector register
#---------------------------------------
type Vreg <: Operand
    num :: Int # register number, or -1 if not yet assigned a register
    ref :: Int # refcount for lives.  0 -> free after last use
    Vreg() = new(-1,0)
end

# Beginning of lifetime
function bol(x::Vreg) 
    @assert x.ref>=0
    x.ref += 1
    x
end

# End of lifetime
function eol(x::Vreg) 
    @assert x.ref>0
    x.ref -= 1
    x
end

# Reset list of available registers, to registers 0:n-1 on machine
# Setting n less than 16 is useful for reducing the number of REX prefixes
function reset(n::Integer)
    global vregAvail 
    vregAvail = Int[i for i=0:n-1] 
    nothing
end

# Get assembly-code string
function asm(x::Vreg) 
    if x.num==-1 
        x.num = pop!(vregAvail)
    end
    "X$(x.num)"
end

function reassign(z::Vreg, x::Vreg) 
    if z.num==-1
        z.num = x.num   
        true
    else
        false
    end
end

reassign(z::Operand, x::Vreg) = false

# Get assembly code string for last use of x, reusing its register for z if z is unassigned.
function lastuse(x::Vreg, z::Operand)
    # @assert x.num>=0
    n = asm(x)
    if x.ref==0
        if !reassign(z,x)
            @assert findfirst(vregAvail,x.num)==0 
            if false
                # FIFO
                insert!(vregAvail,1,x.num)
            else
                # LIFO
                push!(vregAvail,x.num)
            end
        end
        x.num = -1
    end
    n
end

lastuse(x::OperandOrImmediate, z::Operand) = asm(x)

#---------------------------------------
# General-purpose register
#---------------------------------------
type Ireg <: Operand
    asm :: ASCIIString # Name of register
end

asm(x::Ireg) = x.asm

#---------------------------------------
# Location relative to a base register
#---------------------------------------
type Loc <: Mem
    base :: Ireg 	# Base register
    offset :: Int   # Offset
end

#---------------------------------------
# Location relative to a base register
#---------------------------------------
asm(loc::Loc) = "$(loc.offset)($(asm(loc.base)))"

type IndexedLoc <: Mem
    base :: Ireg 
    index :: Ireg
    scale :: Int
end

asm(loc::IndexedLoc) = "($(asm(loc.base)))($(asm(loc.index))*$(loc.scale))"

#---------------------------------------
# Argument to a function.  Go passes all arguments in memory.
#---------------------------------------
type Arg <: Mem
    name :: ASCIIString  # Name of the argument.  Go assembler requires one.
    offset :: Int		 # Offset of argument from FP
end

asm(arg::Arg) = "$(arg.name)+$(arg.offset)(FP)"

const IregOrMemOrInteger = Union{Ireg,Mem,Integer}

#---------------------------------------
# Emitting assembly code
#---------------------------------------
emitln(s::AbstractString...) = println(asmfile, s...)
emitIn(s::AbstractString...) = emitln("\t", s...)   # Emit with indentation

function emit(s::AbstractString, z::Vreg, y::Vreg, x::Integer)
   yasm = lastuse(y,z)
   emitIn(s," ",asm(x), ", ",yasm, ", ",asm(z))
   nothing
end
 
function emit(s::AbstractString, z::Vreg, y::Vreg, x::Vreg) 
   yasm = lastuse(y,z)
   xasm = lastuse(x,z)
   zasm = asm(z)
   if zasm==yasm
       emitIn(s," ",xasm, ", ",zasm)
   else 
       emitIn("MOVAPS ",yasm,", ",zasm)
       emitIn(s," ",xasm, ", ",zasm)
   end
   nothing
end

function emit(s::AbstractString, z::Operand, y::OperandOrImmediate) 
   yasm = lastuse(y,z)
   emitIn(s," ",yasm, ", ",asm(z))
   nothing
end

mov(y::Vreg, x::Vreg) = emit("MOVAPS", y, x)
mov(y::Vreg, x::Mem) = emit("MOVUPS", y, x)
mov(y::Mem, x::Vreg) = emit("MOVUPS", y, x)
mov(y::Ireg, x::Mem) = emit("MOVQ", y, x)
mov(y::Mem, x::Ireg) = emit("MOVQ", y, x)
 
shuf(z::Vreg, x::Vreg, y::Integer) = emit("SHUFPS", z, x, UInt8(y))

function broadcast(z::Vreg, x::Union{Mem,Float32})
    emit("MOVSS",z,x)
    shuf(z,z,0)
end
 
add(z::Vreg, x::Vreg, y::Vreg) = emit("ADDPS", z, x, y)
sub(z::Vreg, x::Vreg, y::Vreg) = emit("SUBPS", z, x, y)
add(z::Vreg, y::Vreg) = emit("ADDPS", z, y)
sub(z::Vreg, y::Vreg) = emit("SUBPS", z, y)
xor(z::Vreg, y::Vreg) = emit("XORPS", z, y)
mul(z::Vreg, x::Vreg, y::Vreg) = emit("MULPS", z, x, y)
sub(z::Ireg, x::IregOrMemOrInteger) = emit("SUBQ", z, x)
add(z::Ireg, x::IregOrMemOrInteger) = emit("ADDQ", z, x)
shl(z::Ireg, x::Integer) = emit("SHLQ", z, x)
cvttss2sq(z::Ireg, x::Vreg) = emit("CVTTSS2SQ", z, x)

#---------------------------------------
# Complex-valued operands
#---------------------------------------
type Cmplx
    re :: Operand
    im :: Operand
    Cmplx() = new(Vreg(),Vreg())
    Cmplx(a::Operand, b::Operand) = new(a,b)
end

function bol(x::Cmplx) 
    bol(x.re)
    bol(x.im)
    x
end

function eol(x::Cmplx) 
    eol(x.re)
    eol(x.im)
    x
end

function mov(y::Cmplx, x::Cmplx)
    mov(y.re, x.re)
    mov(y.im, x.im)
end

function broadcast(y::Cmplx, x::Cmplx)
    broadcast(y.re, x.re)
    broadcast(y.im, x.im)
end

# z = z*x
function mul(z::Cmplx, x::Cmplx)
    a = Vreg()
    b = Vreg()
    bol(x)
    bol(z)
    mul(bol(a),z.im,x.im)
    mul(bol(b),z.re,x.im)
    mul(z.re,z.re,x.re)
    mul(z.im,z.im,x.re)
    sub(z.re,z.re,eol(a))
    add(z.im,z.im,eol(b))
    eol(x)
    eol(z)
end

cmplxLoc(base::Ireg,offset::Integer) = Cmplx(Loc(base,offset), Loc(base,offset+4))

function extractRotate( dst::Ireg, src::Vreg, k::Integer )
   cvttss2sq(dst,src)
   # 0x39 = 00111001
   shuf(src,src,0x39)   # Rotate vector to right
end

function convert4(z::Cmplx, clut::Ireg, dstReg::Ireg, dstOff::Integer) 
    zr = Ireg("AX")
    zi = Ireg("DX")
    tmp = Ireg("BP")
    for k = 0:3
        extractRotate(zr, z.re, k)
        extractRotate(zi, z.im, k)
        shl(zi,7)
        add(zi,zr)
        mov(tmp, IndexedLoc(clut,zi,4))
        mov(Loc(dstReg,dstOff+4*k), tmp)
    end
    # FIXME - need to eol z.re and z.im here
end

#---------------------------------------
# Emit the code
#---------------------------------------
asmfile = open("hft_amd64.s","w")

emitln("#include \"textflag.h\"")

emitln()
reset(16)
emitln("// func accumulateToFeet(z *[2]cvec, u *[2]w13, feet []foot)")
emitln("TEXT ·accumulateToFeet(SB), NOSPLIT, \$$(5*8)")

a = [Cmplx() for i=1:2]
w1 = [Cmplx() for i=1:2]
w3 = [Cmplx() for i=1:2]

fpOff = 0 
const RSize = 16
DI = Ireg("DI")
SI = Ireg("SI")
BX = Ireg("BX")
FP = Ireg("FP")
CX = Ireg("CX")
mov(BX,Arg("z",fpOff)); fpOff+=8
mov(SI,Arg("u",fpOff)); fpOff+=8
for i=1:2
    j = i-1
    mov(bol(a[i]), Cmplx(Loc(BX,0+32j), Loc(BX,16+32j)))
    broadcast(bol(w1[i]), cmplxLoc(SI,0+16j)) 
    broadcast(bol(w3[i]), cmplxLoc(SI,8+16j))
end
mov(DI,Arg("feet",fpOff))
mov(CX,Arg("feet",fpOff+8))
# fpOff+16 has capacity, which is not used here
emitln()

emitln("loop:")

emitIn("// fac, fbc")
fac = Vreg()
fbc = Vreg()
mov(bol(fac),Loc(DI,2*16))
mov(bol(fbc),Loc(DI,3*16))
for i=1:2
    t0 = Vreg()
    t1 = Vreg()
    mul(bol(t0), a[i].re, w1[i].re) 
    mul(bol(t1), a[i].im, w1[i].re) 
    add(fbc,eol(t1))
    add(fac,eol(t0))
end
mov(Loc(DI,2*16), eol(fac))
mov(Loc(DI,3*16), eol(fbc))

emitln()
emitIn("// fad, fbd")
fad = Vreg()
fbd = Vreg()
mov(bol(fad),Loc(DI,4*16))
mov(bol(fbd),Loc(DI,5*16))
for i=1:2
    t0 = Vreg()
    t1 = Vreg()
    mul(bol(t0), a[i].re, w1[i].im) 
    mul(bol(t1), a[i].im, w1[i].im) 
    add(fbd,eol(t1))
    add(fad,eol(t0))
end
mov(Loc(DI,4*16), eol(fad))
mov(Loc(DI,5*16), eol(fbd))

emitln()
emitIn("// fa, fb, rotate")
fa = Vreg()
fb = Vreg()
mov(bol(fa),Loc(DI,0*16))
mov(bol(fb),Loc(DI,1*16))
for i=1:2
    add(fa,a[i].re)
    add(fb,a[i].im)
    mul(a[i], w3[i])
end
mov(Loc(DI,0*16),eol(fa))
mov(Loc(DI,1*16),eol(fb))

emitln()
add(DI,16*6)
sub(CX,1)
emitIn("JG loop")
emitIn("RET")

reset(8)
emitln()
emitln("// func rotate( a []cvec, v [] complex64)")
emitln("TEXT ·rotate(SB), NOSPLIT, \$$(6*8)")
DI = Ireg("DI")
SI = Ireg("SI")
FP = Ireg("FP")
CX = Ireg("CX")
fpOff = 0
mov(DI,Arg("a",fpOff)); fpOff+=8
mov(CX,Arg("a",fpOff)); fpOff+=8
                        fpOff+=8 
mov(SI,Arg("v",fpOff))
emitln("loop:")
v = Cmplx()
broadcast(bol(v), cmplxLoc(SI,0))
a = Cmplx()
aloc = Cmplx(Loc(DI,0), Loc(DI,16))
mov(bol(a),aloc)
mul(a,eol(v))
mov(aloc,eol(a))
add(DI,32) # Advance to next cvec in a
add(SI,8)  # Advance to next cmplx64 in v
sub(CX,1)
emitIn("JG loop")
emitIn("RET")

emitln()
emitln("// func feetToPixel(feet[]foot, clut*[128][128] pixel, row[]pixel)")
emitln("TEXT ·feetToPixel(SB), NOSPLIT, \$$(7*8)")
reset(10)
DI = Ireg("DI") # row
SI = Ireg("SI") # foot
BX = Ireg("BX") # clut
FP = Ireg("FP")
CX = Ireg("CX") # len(foot)
fpOff = 0
mov(SI,Arg("feet",fpOff)); fpOff+=8
mov(CX,Arg("feet",fpOff)); fpOff+=8
fpOff+=8 # Skip cap(feet)
mov(BX,Arg("clut",fpOff)); fpOff+=8
mov(DI,Arg("row",fpOff)); fpOff+=8
zeroVal = Vreg()
xor(zeroVal,bol(zeroVal))
magicVal = Vreg()
broadcast(bol(magicVal),64.5f0)

emitln("loop:")
fad = Vreg()
fac = Vreg()
fbc = Vreg()
fbd = Vreg()
fa = Vreg()
fb = Vreg()
u = Cmplx()
v = Cmplx()
w = Cmplx()
mov(bol(fac),Loc(SI,2*16))
mov(bol(fbc),Loc(SI,3*16))
mov(bol(fad),Loc(SI,4*16))
mov(bol(fbd),Loc(SI,5*16))
emitIn()
emitIn("// left")
sub(bol(u.im),fbc,fad)
add(bol(u.re),fac,fbd)
convert4(u,BX,DI,0)
mov(Loc(SI,2*16),magicVal)
mov(Loc(SI,3*16),magicVal)
emitIn()
emitIn("// middle")
mov(bol(v.re),Loc(SI,0*16))
mov(bol(v.im),Loc(SI,1*16))
convert4(v,BX,DI,16)
mov(Loc(SI,0*16),magicVal)
mov(Loc(SI,1*16),magicVal)
emitIn()
emitIn("// right")
add(bol(w.im),eol(fbc),eol(fad))
sub(bol(w.re),eol(fac),eol(fbd))
convert4(w,BX,DI,32)
mov(Loc(SI,4*16),zeroVal)
mov(Loc(SI,5*16),zeroVal)

emitln()
add(DI,4*12)  # 12 pixels per foot
add(SI,6*16)  # Advance to next foot 
sub(CX,1)
emitIn("JG loop")
emitIn("RET")
close(asmfile)

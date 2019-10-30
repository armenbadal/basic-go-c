
'
' Օրինակ 7
'
SUB s0
    PRINT "I am s0"
END SUB

SUB s1(a)
    PRINT a
END SUB

SUB s2(a, b$)
    PRINT a
    PRINT b$
END SUB

SUB Main
    CALL s0
    CALL s1 32
    CALL s2 777, "Hi"
    'CALL s2 777 "Hi" ' սխալ
END SUB

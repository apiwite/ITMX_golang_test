# ITMX Golang Test!

สวัสดีครับ ผมนายอภิวิชญ์ คะมิ นี่เป็นแบบทดสอบสำหรับการสมัครงานตำแหน่ง **Golang Dev** ในขั้นต่อไปจะเป็นการอธิบายการ setup และ run code ที่ได้ทำการเขียนมา หากผิกพลาดส่วนไหน ขออภัยไว้ ณ ที่นี้ด้วย.


# 1.การสร้าง golang REST App 

สามารถ `go run main.go`  ได้เลย เพื่อรับตัวโปรแกรมหลัก 
server จะรับที่ port : `8080`

## SQlite database
### รูปแบบข้อมูล data ใน SQlite

|ตัวแปร                        	|type                         |
|-------------------------------|-----------------------------|
|`'id'`            				|int|
|`name`            				|string            |
|`age`							|int|

### ข้อมูลใน SQlite database ที่ได้บันทึกไว้
ข้อมูลที่ได้สร้างขึ้นมาเพื่อใช้งานเบื้องต้น
`{
"name" : "april",
"age" : 25
},{
"name" : "jeff",
"age" : 23
}`

## การทดสอบ API ของเส้นต่างๆ
### การทดสอบบน POST MAN

**POST api** สำหรับการเพื่มข้อมูลลง database\
method : `POST`\
path : `http://localhost:8080/customers/add`\
ข้อมูล  json ที่จะเพื่ม\
`{
"name" : "prayut",
"age" : 47
}`\
**ผลลัพธ์**\
"ADD Success"

**GET api** สำหรับการเรียกข้อมูลจาก database ทั้งหมด\
method : `GET`\
path : `http://localhost:8080/customers`\
**ผลลัพธ์**\
`ข้อมูลที่มีอยู่ใน DB`

**PUT api** สำหรับการ update ข้อมูลลง database\
ในที่นี้ เราจะทดสอบด้วยการ update id ที่ 3 ที่ได้เพื่มไปเมื่อข้อที่ผ่านมา\
method : `PUT`\
path : `http://localhost:8080/customers/update/:id`\
 **:id จะแทนด้วย id ที่ต้องการ update ข้อมูล เช่น id:2** `/customers/update/2`\
 **ข้อมูลหรือ id สามารถเรียกดูได้จาก `http://localhost:8080/customers` \
ข้อมูล  json ที่จะ update **สามารถเลือกที่จะ update แค่ name หรือ age เพียงอย่างเดียวก็ได้\
`{
"name" : "pravit", 
"age" : 49
}`\
**ผลลัพธ์**\
สำเร็จ `"status": "updated"`\
ไม่สำเร็จ `"ERROR": " Customer not found"`

**GET api** สำหรับการดึง customer data จาก database\
method : `GET`\
path : `http://localhost:8080/customers/:id`\
 **:id จะแทนด้วย id ที่ต้องการ update ข้อมูล เช่น id:2** `/customers/2`\
 **ข้อมูลหรือ id สามารถเรียกดูได้จาก `http://localhost:8080/customers` \
**ผลลัพธ์**\
สำเร็จ `ข้อมูล customer id นั้นๆ`\
ไม่สำเร็จ `"error": "ID:ที่ค้นหา not found"`

**DELETE api** สำหรับการ delete ข้อมูล database\
method : `DELETE`\
path : `http://localhost:8080/customers/id/:id`\
 **:id จะแทนด้วย id ที่ต้องการ delete ข้อมูล เช่น id:2** `/customers/id/2`\
 **ข้อมูลหรือ id สามารถเรียกดูได้จาก `http://localhost:8080/customers` \
**ผลลัพธ์**\
สำเร็จ `"message": "Customer ID : idที่ลบ deleted successfully"`\
ไม่สำเร็จ `"error": "Customer ID not found"`

**GET api** สำหรับการ create และเพื่มข้อมูลลง database\
method : `GET`\
path : `http://localhost:8080/resetDB`\
**หากจะทำการ resetDB ใหม่ ให้เลิก run server แล้วทำการลบไฟล์  `Customers.db` แล้วทำการรันใหม่ แล้วค่อยเรียกใช้ api**

## Unit Test

สามารถ `go test` เพื่อทำการ run test\
ทุกครั้งที่เริ่ม run test จะทำการสร้าง data base ที่สำหรับ test ใหม่ ชื่อ `Customers_test.db` \
โดยมีข้อมูลเริ่มต้น `{
"name" : "pravit", 
"age" : 49
}`\
หรือสามารถเข้าไปยังไฟล์ `main_test.go` แล้วสามารถเลือกรันแต่ละ test ได้เหมือนกัน

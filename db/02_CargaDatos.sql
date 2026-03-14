/********************ROLES*******************/
INSERT INTO Roles(nombre)
    VALUES 
        ('Admin'),
        ('Client')
    ;

/********************USERS*******************/
INSERT INTO Users (nombre, role_id)
    VALUES 
        ('Cesar', 1),
        ('Daniel', 1),
        ('Kirsten', 1)
    ;

INSERT INTO Users (nombre, role_id)
    VALUES 
        ('Franco', 2),
        ('Dieguito', 2),
        ('Benavides', 2)
    ;

/********************RESTAURANT*******************/
INSERT INTO Restaurant (nombre, admin_id, estado)
    VALUES 
        ('Las pizzas bros', 1, 1),
        ('La cucharacha', 2, 1)
    ;

/********************MENU*******************/    

INSERT INTO Menu(dish_name,price,restaurant_id)
    VALUES
        ('Pizza Napolitana',7500.00,1),
        ('Pizza Pepperoni', 8200.00, 1),
        ('Lasagna', 6900.00, 1),
        ('Pasta Alfredo', 6400.00, 1)
    ;

INSERT INTO Menu(dish_name,price,restaurant_id)
    VALUES
        ('Carne mechada',7500.00,2),
        ('Casado', 8200.00, 2),
        ('Penne', 6900.00, 2),
        ('Lomito', 6400.00, 2)
    ;

/********************TABLES*******************/  
INSERT INTO Tables(table_number, estado, restaurant_id)
    VALUES 
        (1, 1, 1),
        (2, 1, 1),
        (3, 1, 1),
        (4, 1, 1)
    ;

INSERT INTO Tables(table_number, estado, restaurant_id)
    VALUES 
        (1, 1, 2),
        (2, 1, 2),
        (3, 1, 2),
        (4, 1, 2)
    ;

/********************ORDERS*******************/  
INSERT INTO Orders(table_id, client_id, orders_type, restaurant_id)
    VALUES
        (1, 4, 'Dine in', 1),
        (2, 5, 'Dine in', 1),
        (2, 6, 'Dine in', 1)
    ;

INSERT INTO Orders(client_id, orders_type, restaurant_id)
    VALUES
        (4, 'To go', 2),
        (5, 'To go', 2),
        (6, 'To go', 2)
    ;

/********************ORDER_ITEMS*******************/

INSERT INTO Orders_Items (orders_id, menu_id, quantity, price)
    VALUES
        (1, 1, 2, 7500.00),
        (1, 2, 1, 8200.00),
        (2, 3, 1, 6900.00),
        (3, 4, 3, 6400.00),
        (4, 5, 2, 7500.00),
        (5, 6, 1, 8200.00),
        (6, 7, 2, 6900.00)
    ;
/********************RESERVATIONS*******************/
INSERT INTO Reservation (table_id, client_id, fecha, estado)
    VALUES
        (1, 4, '2024-06-10 19:00:00', 1),
        (2, 5, '2024-06-11 20:00:00', 1),
        (3, 6, '2024-06-12 21:00:00', 1)
    ;
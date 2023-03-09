DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`(
    `user_id` int(11) NOT NULL AUTO_INCREMENT,
    `username` varchar(100) NOT NULL,
    `customer_id` int(11) NOT NULL,
    `role` varchar(10) NOT NULL,
    `password` varchar(100) NOT NULL,
    PRIMARY KEY (`user_id`)
)   ENGINE=InnoDB AUTO_INCREMENT=2106 DEFAULT CHARSET=latin1;

DROP TABLE IF EXISTS `customers`;
CREATE TABLE `customers`(
    `customer_id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(100) NOT NULL,
    `date_of_birth` date NOT NULL,
    `city` varchar(100) NOT NULL,
    `zipcode` varchar(10) NOT NULL,
    `status` tinyint(1) NOT NULL DEFAULT '1',
    PRIMARY KEY (`customer_id`)
)   ENGINE=InnoDB AUTO_INCREMENT=2006 DEFAULT CHARSET=latin1;

DROP TABLE IF EXISTS `accounts`;
CREATE TABLE `accounts`(
    `account_id` int(11) NOT NULL AUTO_INCREMENT,
    `customer_id` int(11) NOT NULL,
    `opening_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `account_type` varchar(10) NOT NULL,
    `pin` varchar(10) NOT NULL,
    `status` tinyint(1) NOT NULL DEFAULT '1',
    PRIMARY KEY (`account_id`),
    KEY `accounts_FK` (`customer_id`),
    CONSTRAINT `accounts_FK` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`customer_id`)
)   ENGINE=InnoDB AUTO_INCREMENT=95476 DEFAULT CHARSET=latin1;


DROP TABLE IF EXISTS `transactions`;
CREATE TABLE `transactions`(
    `transaction_id` int(11) NOT NULL AUTO_INCREMENT,
    `account_id` int(11) NOT NULL,
    `amount` int(11) NOT NULL,
    `transaction_type` varchar(10) NOT NULL,
    `transaction_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`transaction_id`),
    KEY `transactions_FK` (`account_id`),
    CONSTRAINT `transactions_FK` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`account_id`)
)   ENGINE=InnoDB DEFAULT CHARSET=latin1;

SELECT username, u.customer_id, u.role, GROUP_CONCAT(a.account_id) as account_numbers
FROM users u LEFT JOIN accounts a ON a.customer_id = u.customer_id 
WHERE username="admin" and password="abc123"
GROUP BY username, u.customer_id,u.role;